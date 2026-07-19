// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all environment-driven settings for the API server.
type Config struct {
	Env                string
	HTTPPort           string
	DatabaseURL        string
	CORSOrigins        []string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	AuthServiceURL     string
	InternalServiceKey string
	DefaultBrandSlug   string
}

func Load() Config {
	return Config{
		Env:                getEnv("APP_ENV", "development"),
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://didar:didar@localhost:5432/didar?sslmode=disable"),
		CORSOrigins:        splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173")),
		JWTSecret:          getEnv("JWT_SECRET", "dev-insecure-secret-change-me"),
		AccessTokenTTL:     getDurationMinutes("ACCESS_TOKEN_TTL_MINUTES", 30),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		InternalServiceKey: getEnv("INTERNAL_SERVICE_KEY", "dev-internal-key-change-me"),
		DefaultBrandSlug:   getEnv("DEFAULT_BRAND_SLUG", "didar"),
	}
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
