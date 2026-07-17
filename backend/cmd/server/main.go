// Command server boots the Didar API: it loads config, connects to
// Postgres/PostGIS, self-migrates the schema, and serves the HTTP API.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/didar-platform/backend/internal/auth"
	"github.com/didar-platform/backend/internal/booking"
	"github.com/didar-platform/backend/internal/brand"
	"github.com/didar-platform/backend/internal/config"
	"github.com/didar-platform/backend/internal/dbx"
	"github.com/didar-platform/backend/internal/experience"
	"github.com/didar-platform/backend/internal/httpserver"
	"github.com/didar-platform/backend/internal/otp"
	"github.com/didar-platform/backend/internal/user"
	"github.com/didar-platform/backend/migrations"
)

func main() {
	_ = godotenv.Load()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := dbx.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := dbx.Migrate(ctx, pool, migrations.FS); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	var sender otp.SMSSender
	if cfg.OTPMockMode {
		slog.Warn("OTP_MOCK_MODE is on: OTP codes will be logged instead of sent via SMS")
		sender = otp.MockSender{}
	} else {
		sender = otp.NewKavenegarSender(cfg.KavenegarAPIKey, cfg.KavenegarTemplate, cfg.KavenegarSenderLine)
	}

	userRepo := user.NewRepository(pool)
	brandRepo := brand.NewRepository(pool)
	experienceRepo := experience.NewRepository(pool)
	bookingRepo := booking.NewRepository(pool)

	otpService := otp.NewService(pool, sender, cfg.OTPTTL, cfg.OTPCodeLen)
	tokenIssuer := auth.NewTokenIssuer(cfg.JWTSecret, cfg.AccessTokenTTL)
	sessionRepo := auth.NewSessionRepository(pool, cfg.RefreshTokenTTL)

	authHandler := auth.NewHandler(userRepo, otpService, tokenIssuer, sessionRepo, brandRepo, cfg.DefaultBrandSlug)
	brandHandler := brand.NewHandler(brandRepo, cfg.DefaultBrandSlug)
	experienceHandler := experience.NewHandler(experienceRepo)
	experienceHostHandler := experience.NewHostHandler(experienceRepo, func() (string, error) {
		b, err := brandRepo.GetBySlug(context.Background(), cfg.DefaultBrandSlug)
		if err != nil {
			return "", err
		}
		return b.ID, nil
	})
	bookingHandler := booking.NewHandler(bookingRepo)

	router := httpserver.NewRouter(httpserver.Deps{
		CORSOrigins:       cfg.CORSOrigins,
		TokenIssuer:       tokenIssuer,
		AuthHandler:       authHandler,
		BrandHandler:      brandHandler,
		ExperienceHandler: experienceHandler,
		ExperienceHost:    experienceHostHandler,
		BookingHandler:    bookingHandler,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server listening", "port", cfg.HTTPPort, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
