// Command server boots the Didar auth microservice: OTP login, JWT issuance,
// user identity, and session management.
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

	"github.com/didar-platform/pkg/authkit"
	"github.com/didar-platform/services/auth/internal/auth"
	"github.com/didar-platform/services/auth/internal/config"
	"github.com/didar-platform/services/auth/internal/dbx"
	"github.com/didar-platform/services/auth/internal/httpserver"
	"github.com/didar-platform/services/auth/internal/otp"
	"github.com/didar-platform/services/auth/internal/user"
	"github.com/didar-platform/services/auth/migrations"
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

	var sender otp.OTPSender
	if cfg.OTPMockMode {
		slog.Warn("OTP_MOCK_MODE is on: OTP codes will be logged instead of sent")
		sender = otp.MockSender{}
	} else if cfg.OTPChannel == "sms" {
		slog.Info("OTP channel: sms (Kavenegar)")
		sender = otp.NewKavenegarSender(cfg.KavenegarAPIKey, cfg.KavenegarTemplate, cfg.KavenegarSenderLine)
	} else {
		slog.Info("OTP channel: email (SMTP)")
		sender = otp.NewEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom)
	}

	userRepo := user.NewRepository(pool)
	otpService := otp.NewService(pool, sender, cfg.OTPTTL, cfg.OTPCodeLen)
	tokenIssuer := authkit.NewTokenIssuer(cfg.JWTSecret, cfg.AccessTokenTTL)
	sessionRepo := auth.NewSessionRepository(pool, cfg.RefreshTokenTTL)
	authHandler := auth.NewHandler(userRepo, otpService, tokenIssuer, sessionRepo)

	router := httpserver.NewRouter(httpserver.Deps{
		CORSOrigins:        cfg.CORSOrigins,
		TokenIssuer:        tokenIssuer,
		InternalServiceKey: cfg.InternalServiceKey,
		AuthHandler:        authHandler,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("auth service listening", "port", cfg.HTTPPort, "env", cfg.Env)
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
