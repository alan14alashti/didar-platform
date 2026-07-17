package experience

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/didar-platform/backend/internal/auth"
	"github.com/didar-platform/backend/internal/platform"
)

type HostHandler struct {
	repo        *Repository
	defaultBrandID func() (string, error)
}

// NewHostHandler takes a lazily-resolved default brand id getter so the
// host doesn't need to know/send a brand_id when creating an experience —
// it is inferred from the currently resolved white-label brand.
func NewHostHandler(repo *Repository, defaultBrandID func() (string, error)) *HostHandler {
	return &HostHandler{repo: repo, defaultBrandID: defaultBrandID}
}

func hostID(r *http.Request) (string, bool) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		return "", false
	}
	return claims.UserID, true
}

func (h *HostHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	items, err := h.repo.ListByHost(r.Context(), uid)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *HostHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	e, err := h.repo.GetByIDForHost(r.Context(), chi.URLParam(r, "id"), uid)
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, e)
}

type experienceBody struct {
	CategoryID      *string  `json:"category_id"`
	Title           string   `json:"title"`
	Summary         *string  `json:"summary"`
	Description     *string  `json:"description"`
	PriceAmount     float64  `json:"price_amount"`
	DurationMinutes int      `json:"duration_minutes"`
	MaxGuests       int      `json:"max_guests"`
	MinGuests       int      `json:"min_guests"`
	City            string   `json:"city"`
	Address         *string  `json:"address"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	Status          string   `json:"status"`
}

func (b experienceBody) validate() error {
	if b.Title == "" {
		return platform.ErrValidation("title is required")
	}
	if b.City == "" {
		return platform.ErrValidation("city is required")
	}
	if b.PriceAmount <= 0 {
		return platform.ErrValidation("price_amount must be greater than zero")
	}
	if b.Latitude == 0 && b.Longitude == 0 {
		return platform.ErrValidation("a valid location (latitude/longitude) is required")
	}
	if b.Status != "" && b.Status != "draft" && b.Status != "published" && b.Status != "archived" {
		return platform.ErrValidation("status must be draft, published or archived")
	}
	return nil
}

func (b experienceBody) toInput() UpsertExperienceInput {
	durationMinutes := b.DurationMinutes
	if durationMinutes <= 0 {
		durationMinutes = 120
	}
	maxGuests := b.MaxGuests
	if maxGuests <= 0 {
		maxGuests = 6
	}
	minGuests := b.MinGuests
	if minGuests <= 0 {
		minGuests = 1
	}
	status := b.Status
	if status == "" {
		status = "draft"
	}
	return UpsertExperienceInput{
		CategoryID:      b.CategoryID,
		Title:           b.Title,
		Summary:         b.Summary,
		Description:     b.Description,
		PriceAmount:     b.PriceAmount,
		DurationMinutes: durationMinutes,
		MaxGuests:       maxGuests,
		MinGuests:       minGuests,
		City:            b.City,
		Address:         b.Address,
		Latitude:        b.Latitude,
		Longitude:       b.Longitude,
		Status:          status,
	}
}

func (h *HostHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body experienceBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if err := body.validate(); err != nil {
		platform.WriteError(w, err)
		return
	}

	brandID, err := h.defaultBrandID()
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	e, err := h.repo.Create(r.Context(), brandID, uid, body.toInput())
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusCreated, e)
}

func (h *HostHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body experienceBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if err := body.validate(); err != nil {
		platform.WriteError(w, err)
		return
	}

	e, err := h.repo.Update(r.Context(), chi.URLParam(r, "id"), uid, body.toInput())
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, e)
}

func (h *HostHandler) Archive(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	if err := h.repo.Archive(r.Context(), chi.URLParam(r, "id"), uid); err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

type addPhotoBody struct {
	URL     string  `json:"url"`
	Alt     *string `json:"alt"`
	IsCover bool    `json:"is_cover"`
}

func (h *HostHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body addPhotoBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if body.URL == "" {
		platform.WriteError(w, platform.ErrValidation("url is required"))
		return
	}

	p, err := h.repo.AddPhoto(r.Context(), chi.URLParam(r, "id"), uid, body.URL, body.Alt, body.IsCover)
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusCreated, p)
}

func (h *HostHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	err := h.repo.DeletePhoto(r.Context(), chi.URLParam(r, "id"), uid, chi.URLParam(r, "photoId"))
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

type addSlotBody struct {
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Capacity int       `json:"capacity"`
}

func (h *HostHandler) AddSlot(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body addSlotBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if body.Capacity <= 0 || !body.EndAt.After(body.StartAt) {
		platform.WriteError(w, platform.ErrValidation("capacity must be positive and end_at must be after start_at"))
		return
	}

	slot, err := h.repo.AddSlot(r.Context(), chi.URLParam(r, "id"), uid, Slot{StartAt: body.StartAt, EndAt: body.EndAt, Capacity: body.Capacity})
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusCreated, slot)
}

func (h *HostHandler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	uid, ok := hostID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	err := h.repo.DeleteSlot(r.Context(), chi.URLParam(r, "id"), uid, chi.URLParam(r, "slotId"))
	if err != nil {
		writeRepoErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func writeRepoErr(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		platform.WriteError(w, platform.ErrNotFound)
		return
	}
	platform.WriteError(w, err)
}
