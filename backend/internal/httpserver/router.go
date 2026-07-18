// Package httpserver assembles all feature handlers into a single chi
// router with shared middleware (CORS, logging, panic recovery, auth).
package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/didar-platform/backend/internal/auth"
	"github.com/didar-platform/backend/internal/booking"
	"github.com/didar-platform/backend/internal/brand"
	"github.com/didar-platform/backend/internal/experience"
	"github.com/didar-platform/backend/internal/platform"
)

type Deps struct {
	CORSOrigins []string
	TokenIssuer *auth.TokenIssuer

	AuthHandler       *auth.Handler
	BrandHandler       *brand.Handler
	ExperienceHandler  *experience.Handler
	ExperienceHost     *experience.HostHandler
	BookingHandler     *booking.Handler
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Bypass-Tunnel-Reminder"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		platform.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	requireAuth := auth.RequireAuth(d.TokenIssuer)

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/brand", d.BrandHandler.Current)
		api.Get("/categories", d.ExperienceHandler.Categories)
		api.Get("/experiences", d.ExperienceHandler.List)
		api.Get("/experiences/{slug}", d.ExperienceHandler.GetBySlug)

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
			g.Post("/me/become-host", d.AuthHandler.BecomeHost)

			g.Post("/bookings", d.BookingHandler.Create)
			g.Get("/me/bookings", d.BookingHandler.ListMine)
			g.Post("/me/bookings/{id}/cancel", d.BookingHandler.CancelMine)

			g.Route("/host", func(h chi.Router) {
				h.Use(auth.RequireHost)

				h.Get("/experiences", d.ExperienceHost.ListMine)
				h.Post("/experiences", d.ExperienceHost.Create)
				h.Get("/experiences/{id}", d.ExperienceHost.Get)
				h.Patch("/experiences/{id}", d.ExperienceHost.Update)
				h.Delete("/experiences/{id}", d.ExperienceHost.Archive)
				h.Post("/experiences/{id}/photos", d.ExperienceHost.AddPhoto)
				h.Delete("/experiences/{id}/photos/{photoId}", d.ExperienceHost.DeletePhoto)
				h.Post("/experiences/{id}/slots", d.ExperienceHost.AddSlot)
				h.Delete("/experiences/{id}/slots/{slotId}", d.ExperienceHost.DeleteSlot)

				h.Get("/bookings", d.BookingHandler.ListForHost)
				h.Patch("/bookings/{id}", d.BookingHandler.UpdateStatusForHost)
			})
		})
	})

	return r
}
