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

	"github.com/didar-platform/backend/internal/authclient"
	"github.com/didar-platform/backend/internal/booking"
	"github.com/didar-platform/backend/internal/brand"
	"github.com/didar-platform/backend/internal/config"
	"github.com/didar-platform/backend/internal/dbx"
	"github.com/didar-platform/backend/internal/experience"
	"github.com/didar-platform/backend/internal/host"
	"github.com/didar-platform/backend/internal/httpserver"
	"github.com/didar-platform/backend/internal/user"
	"github.com/didar-platform/backend/migrations"
	"github.com/didar-platform/pkg/authkit"
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

	userRepo := user.NewRepository(pool)
	brandRepo := brand.NewRepository(pool)
	experienceRepo := experience.NewRepository(pool)
	bookingRepo := booking.NewRepository(pool)

	tokenIssuer := authkit.NewTokenIssuer(cfg.JWTSecret, cfg.AccessTokenTTL)
	authClient := authclient.New(cfg.AuthServiceURL, cfg.InternalServiceKey)

	hostHandler := host.NewHandler(userRepo, brandRepo, authClient, cfg.DefaultBrandSlug)
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
		HostHandler:       hostHandler,
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
