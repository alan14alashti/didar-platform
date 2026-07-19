package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/didar-platform/pkg/authkit"
	"github.com/didar-platform/pkg/platform"
	"github.com/didar-platform/services/auth/internal/auth"
)

type Deps struct {
	CORSOrigins        []string
	TokenIssuer        *authkit.TokenIssuer
	InternalServiceKey string
	AuthHandler        *auth.Handler
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(20 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Service-Key", "Bypass-Tunnel-Reminder"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		platform.JSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "auth"})
	})

	requireAuth := authkit.RequireAuth(d.TokenIssuer)

	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/auth", func(a chi.Router) {
			a.Post("/otp/request", d.AuthHandler.RequestOTP)
			a.Post("/otp/verify", d.AuthHandler.VerifyOTP)
			a.Post("/refresh", d.AuthHandler.Refresh)
			a.Post("/logout", d.AuthHandler.Logout)
		})

		api.Group(func(g chi.Router) {
			g.Use(requireAuth)
			g.Get("/me", d.AuthHandler.Me)
			g.Patch("/me", d.AuthHandler.UpdateMe)
		})
	})

	r.Route("/internal", func(internal chi.Router) {
		internal.Use(requireServiceKey(d.InternalServiceKey))
		internal.Post("/users/{id}/enable-host", func(w http.ResponseWriter, r *http.Request) {
			d.AuthHandler.EnableHost(w, r, chi.URLParam(r, "id"))
		})
	})

	return r
}

func requireServiceKey(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expected == "" || r.Header.Get("X-Service-Key") != expected {
				platform.WriteError(w, platform.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
