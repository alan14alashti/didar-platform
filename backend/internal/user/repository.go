package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const userColumns = `id::text, phone, full_name, avatar_url, locale, is_host, is_admin, status, last_login_at, created_at, updated_at`

func scanUser(row pgx.Row) (*User, error) {
	var u User
	if err := row.Scan(&u.ID, &u.Phone, &u.FullName, &u.AvatarURL, &u.Locale, &u.IsHost, &u.IsAdmin, &u.Status, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE phone = $1`, phone)
	return scanUser(row)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *Repository) CreateWithPhone(ctx context.Context, phone string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO users (phone, last_login_at) VALUES ($1, now())
		RETURNING `+userColumns, phone)
	return scanUser(row)
}

func (r *Repository) TouchLastLogin(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = now() WHERE id = $1`, id)
	return err
}

func (r *Repository) UpdateProfile(ctx context.Context, id string, fullName *string, locale *string, avatarURL *string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE users SET
			full_name = COALESCE($2, full_name),
			locale = COALESCE($3, locale),
			avatar_url = COALESCE($4, avatar_url),
			updated_at = now()
		WHERE id = $1
		RETURNING `+userColumns, id, fullName, locale, avatarURL)
	return scanUser(row)
}

func (r *Repository) BecomeHost(ctx context.Context, id string, brandID string, displayName *string, bio *string, city *string) (*User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE users SET is_host = true, updated_at = now() WHERE id = $1`, id); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO host_profiles (user_id, brand_id, display_name, bio, city)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			bio = EXCLUDED.bio,
			city = EXCLUDED.city,
			updated_at = now()
	`, id, brandID, displayName, bio, city); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) GetHostProfile(ctx context.Context, userID string) (*HostProfile, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT user_id::text, display_name, bio, city, is_verified, rating_avg, rating_count, created_at
		FROM host_profiles WHERE user_id = $1
	`, userID)
	var hp HostProfile
	if err := row.Scan(&hp.UserID, &hp.DisplayName, &hp.Bio, &hp.City, &hp.IsVerified, &hp.RatingAvg, &hp.RatingCount, &hp.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &hp, nil
}
