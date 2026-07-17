package experience

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/didar-platform/backend/internal/platform"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.ListCategories(r.Context())
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"items": cats})
}

func parseFloatPtr(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

// List handles GET /experiences with optional text/category/city filters
// and optional lat/lng/radius_km for proximity search & sorting.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	radius, _ := strconv.ParseFloat(q.Get("radius_km"), 64)

	params := SearchParams{
		Query:        q.Get("query"),
		CategorySlug: q.Get("category"),
		City:         q.Get("city"),
		Lat:          parseFloatPtr(q.Get("lat")),
		Lng:          parseFloatPtr(q.Get("lng")),
		RadiusKM:     radius,
		Page:         page,
		PageSize:     pageSize,
	}

	result, err := h.repo.Search(r.Context(), params)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	e, err := h.repo.GetBySlugPublished(r.Context(), slug)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			platform.WriteError(w, platform.ErrNotFound)
			return
		}
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, e)
}
