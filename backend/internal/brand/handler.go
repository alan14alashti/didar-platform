package brand

import (
	"net/http"
	"strings"

	"github.com/didar-platform/backend/internal/platform"
)

type Handler struct {
	repo        *Repository
	defaultSlug string
}

func NewHandler(repo *Repository, defaultSlug string) *Handler {
	return &Handler{repo: repo, defaultSlug: defaultSlug}
}

// Current resolves & returns the brand for the requesting host. The
// frontend calls this once on boot to theme itself (colors, name, logo,
// locales, currency) — this is the crux of the white-label design.
func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	b, err := h.repo.Resolve(r.Context(), host, h.defaultSlug)
	if err != nil {
		platform.WriteError(w, platform.ErrNotFound)
		return
	}
	platform.JSON(w, http.StatusOK, b)
}
