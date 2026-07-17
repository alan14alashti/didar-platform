package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found or revoked")

type SessionRepository struct {
	pool *pgxpool.Pool
	ttl  time.Duration
}

func NewSessionRepository(pool *pgxpool.Pool, ttl time.Duration) *SessionRepository {
	return &SessionRepository{pool: pool, ttl: ttl}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generateOpaqueToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// Create issues a new refresh token, stores its hash, and returns the
// plaintext token to hand back to the client.
func (r *SessionRepository) Create(ctx context.Context, userID, userAgent, ipAddress string) (string, error) {
	token, err := generateOpaqueToken()
	if err != nil {
		return "", err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO sessions (user_id, refresh_token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, hashToken(token), userAgent, ipAddress, time.Now().Add(r.ttl))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Rotate validates the given refresh token, revokes it, and issues a new one.
func (r *SessionRepository) Rotate(ctx context.Context, token, userAgent, ipAddress string) (userID string, newToken string, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)

	var sessionID string
	var expiresAt time.Time
	var revokedAt *time.Time
	err = tx.QueryRow(ctx, `
		SELECT id::text, user_id::text, expires_at, revoked_at FROM sessions WHERE refresh_token_hash = $1
	`, hashToken(token)).Scan(&sessionID, &userID, &expiresAt, &revokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrSessionNotFound
		}
		return "", "", err
	}
	if revokedAt != nil || time.Now().After(expiresAt) {
		return "", "", ErrSessionNotFound
	}

	if _, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE id = $1`, sessionID); err != nil {
		return "", "", err
	}

	newToken, err = generateOpaqueToken()
	if err != nil {
		return "", "", err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO sessions (user_id, refresh_token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, hashToken(newToken), userAgent, ipAddress, time.Now().Add(r.ttl)); err != nil {
		return "", "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", err
	}

	return userID, newToken, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE refresh_token_hash = $1`, hashToken(token))
	return err
}
