package otp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrTooManyRequests  = errors.New("otp requested too recently")
	ErrExpiredOrUnknown = errors.New("otp code expired or not found")
	ErrIncorrectCode    = errors.New("otp code incorrect")
	ErrTooManyAttempts  = errors.New("too many incorrect attempts")
)

const PurposeLogin = "login"

const resendCooldown = 30 * time.Second

// OTPSender delivers a one-time code to a recipient (email or phone).
type OTPSender interface {
	SendOTP(ctx context.Context, recipient, code string) error
}

type Service struct {
	pool    *pgxpool.Pool
	sender  OTPSender
	ttl     time.Duration
	codeLen int
}

func NewService(pool *pgxpool.Pool, sender OTPSender, ttl time.Duration, codeLen int) *Service {
	if codeLen <= 0 {
		codeLen = 5
	}
	return &Service{pool: pool, sender: sender, ttl: ttl, codeLen: codeLen}
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func generateNumericCode(length int) (string, error) {
	digits := make([]byte, length)
	max := big.NewInt(10)
	for i := range digits {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		digits[i] = byte('0' + n.Int64())
	}
	return string(digits), nil
}

func (s *Service) Request(ctx context.Context, recipient, purpose, ipAddress string) error {
	var lastCreated time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT created_at FROM otp_codes
		WHERE recipient = $1 AND purpose = $2
		ORDER BY created_at DESC LIMIT 1
	`, recipient, purpose).Scan(&lastCreated)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check last otp: %w", err)
	}
	if err == nil && time.Since(lastCreated) < resendCooldown {
		return ErrTooManyRequests
	}

	code, err := generateNumericCode(s.codeLen)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO otp_codes (recipient, purpose, code_hash, expires_at, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`, recipient, purpose, hashCode(code), time.Now().Add(s.ttl), ipAddress)
	if err != nil {
		return fmt.Errorf("store otp: %w", err)
	}

	if err := s.sender.SendOTP(ctx, recipient, code); err != nil {
		return fmt.Errorf("send otp: %w", err)
	}

	return nil
}

func (s *Service) Verify(ctx context.Context, recipient, purpose, code string) error {
	var (
		id           string
		codeHash     string
		attemptCount int
		maxAttempts  int
		expiresAt    time.Time
		consumedAt   *time.Time
	)

	err := s.pool.QueryRow(ctx, `
		SELECT id::text, code_hash, attempt_count, max_attempts, expires_at, consumed_at
		FROM otp_codes
		WHERE recipient = $1 AND purpose = $2
		ORDER BY created_at DESC LIMIT 1
	`, recipient, purpose).Scan(&id, &codeHash, &attemptCount, &maxAttempts, &expiresAt, &consumedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrExpiredOrUnknown
		}
		return fmt.Errorf("load otp: %w", err)
	}

	if consumedAt != nil || time.Now().After(expiresAt) {
		return ErrExpiredOrUnknown
	}
	if attemptCount >= maxAttempts {
		return ErrTooManyAttempts
	}

	if hashCode(code) != codeHash {
		if _, err := s.pool.Exec(ctx, `UPDATE otp_codes SET attempt_count = attempt_count + 1 WHERE id = $1`, id); err != nil {
			return fmt.Errorf("increment attempt count: %w", err)
		}
		return ErrIncorrectCode
	}

	if _, err := s.pool.Exec(ctx, `UPDATE otp_codes SET consumed_at = now() WHERE id = $1`, id); err != nil {
		return fmt.Errorf("consume otp: %w", err)
	}

	return nil
}
