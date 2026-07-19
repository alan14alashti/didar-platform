// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all environment-driven settings for the auth service.
type Config struct {
	Env                string
	HTTPPort           string
	DatabaseURL        string
	CORSOrigins        []string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	InternalServiceKey string

	// OTPChannel is "email" (current) or "sms" (phone/Kavenegar, for later restore).
	OTPChannel string
	OTPMockMode bool
	OTPTTL      time.Duration
	OTPCodeLen  int

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Kept for restoring phone SMS auth later.
	KavenegarAPIKey     string
	KavenegarTemplate   string
	KavenegarSenderLine string
}

func Load() Config {
	cfg := Config{
		Env:                getEnv("APP_ENV", "development"),
		HTTPPort:           getEnv("HTTP_PORT", "8081"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://didar:didar@localhost:5432/didar_auth?sslmode=disable"),
		CORSOrigins:        splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173")),
		JWTSecret:          getEnv("JWT_SECRET", "dev-insecure-secret-change-me"),
		AccessTokenTTL:     getDurationMinutes("ACCESS_TOKEN_TTL_MINUTES", 30),
		RefreshTokenTTL:    getDurationMinutes("REFRESH_TOKEN_TTL_MINUTES", 60*24*30),
		InternalServiceKey: getEnv("INTERNAL_SERVICE_KEY", "dev-internal-key-change-me"),
		OTPChannel:         strings.ToLower(getEnv("OTP_CHANNEL", "email")),
		OTPTTL:             getDurationMinutes("OTP_TTL_MINUTES", 2),
		OTPCodeLen:         getInt("OTP_CODE_LENGTH", 5),
		SMTPHost:           getEnv("SMTP_HOST", ""),
		SMTPPort:           getInt("SMTP_PORT", 587),
		SMTPUsername:       getEnv("SMTP_USERNAME", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:           getEnv("SMTP_FROM", "noreply@didar.travel"),
		KavenegarAPIKey:    getEnv("KAVENEGAR_API_KEY", ""),
		KavenegarTemplate:  getEnv("KAVENEGAR_OTP_TEMPLATE", "didar-otp"),
		KavenegarSenderLine: getEnv("KAVENEGAR_SENDER_LINE", ""),
	}

	defaultMock := true
	switch cfg.OTPChannel {
	case "sms":
		defaultMock = cfg.KavenegarAPIKey == ""
	default:
		defaultMock = cfg.SMTPHost == ""
	}
	cfg.OTPMockMode = getBool("OTP_MOCK_MODE", defaultMock)
	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getDurationMinutes(key string, fallbackMinutes int) time.Duration {
	return time.Duration(getInt(key, fallbackMinutes)) * time.Minute
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
