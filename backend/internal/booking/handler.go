package booking

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/didar-platform/backend/internal/auth"
	"github.com/didar-platform/backend/internal/platform"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func callerID(r *http.Request) (string, bool) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		return "", false
	}
	return claims.UserID, true
}

type createBookingBody struct {
	SlotID      string  `json:"slot_id"`
	GuestsCount int     `json:"guests_count"`
	GuestNote   *string `json:"guest_note"`
}

// Create handles a guest's booking request for an experience time slot.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := callerID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body createBookingBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if body.SlotID == "" || body.GuestsCount <= 0 {
		platform.WriteError(w, platform.ErrValidation("slot_id and a positive guests_count are required"))
		return
	}

	b, err := h.repo.Create(r.Context(), uid, body.SlotID, body.GuestsCount, body.GuestNote)
	if err != nil {
		writeBookingErr(w, err)
		return
	}
	platform.JSON(w, http.StatusCreated, b)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	uid, ok := callerID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	items, err := h.repo.ListByGuest(r.Context(), uid)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) CancelMine(w http.ResponseWriter, r *http.Request) {
	uid, ok := callerID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	if err := h.repo.CancelByGuest(r.Context(), chi.URLParam(r, "id"), uid); err != nil {
		writeBookingErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) ListForHost(w http.ResponseWriter, r *http.Request) {
	uid, ok := callerID(r)
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

type updateStatusBody struct {
	Status   string  `json:"status"`
	HostNote *string `json:"host_note"`
}

func (h *Handler) UpdateStatusForHost(w http.ResponseWriter, r *http.Request) {
	uid, ok := callerID(r)
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body updateStatusBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if body.Status != StatusConfirmed && body.Status != StatusCancelled && body.Status != StatusCompleted {
		platform.WriteError(w, platform.ErrValidation("status must be confirmed, cancelled or completed"))
		return
	}

	b, err := h.repo.UpdateStatusByHost(r.Context(), chi.URLParam(r, "id"), uid, body.Status, body.HostNote)
	if err != nil {
		writeBookingErr(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, b)
}

func writeBookingErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrSlotNotFound):
		platform.WriteError(w, platform.ErrNotFound)
	case errors.Is(err, ErrSlotFull):
		platform.WriteError(w, platform.NewAPIError(http.StatusConflict, "slot_full", "This time slot is fully booked"))
	case errors.Is(err, ErrSlotNotOpen):
		platform.WriteError(w, platform.NewAPIError(http.StatusConflict, "slot_closed", "This time slot is no longer open for booking"))
	case errors.Is(err, ErrInvalidStatus):
		platform.WriteError(w, platform.NewAPIError(http.StatusConflict, "invalid_status_transition", "This booking cannot be updated to that status"))
	default:
		platform.WriteError(w, err)
	}
}
