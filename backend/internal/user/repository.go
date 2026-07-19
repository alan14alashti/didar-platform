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

func (r *Repository) UpsertHostProfile(ctx context.Context, userID, brandID string, displayName *string, bio *string, city *string) (*HostProfile, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO host_profiles (user_id, brand_id, display_name, bio, city)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			bio = EXCLUDED.bio,
			city = EXCLUDED.city,
			updated_at = now()
		RETURNING user_id::text, display_name, bio, city, is_verified, rating_avg, rating_count, created_at
	`, userID, brandID, displayName, bio, city)
	var hp HostProfile
	if err := row.Scan(&hp.UserID, &hp.DisplayName, &hp.Bio, &hp.City, &hp.IsVerified, &hp.RatingAvg, &hp.RatingCount, &hp.CreatedAt); err != nil {
		return nil, err
	}
	return &hp, nil
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
