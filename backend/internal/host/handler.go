package host

import (
	"net/http"

	"github.com/didar-platform/backend/internal/authclient"
	"github.com/didar-platform/backend/internal/brand"
	"github.com/didar-platform/backend/internal/platform"
	"github.com/didar-platform/backend/internal/user"
	"github.com/didar-platform/pkg/authkit"
)

type Handler struct {
	users       *user.Repository
	brands      *brand.Repository
	authClient  *authclient.Client
	defaultSlug string
}

func NewHandler(users *user.Repository, brands *brand.Repository, authClient *authclient.Client, defaultBrandSlug string) *Handler {
	return &Handler{users: users, brands: brands, authClient: authClient, defaultSlug: defaultBrandSlug}
}

type becomeHostBody struct {
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	City        *string `json:"city"`
}

// BecomeHost creates the host profile in the main API database and asks the
// auth service to flip is_host and issue fresh tokens.
func (h *Handler) BecomeHost(w http.ResponseWriter, r *http.Request) {
	claims, ok := authkit.FromContext(r.Context())
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}

	var body becomeHostBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}

	b, err := h.brands.GetBySlug(r.Context(), h.defaultSlug)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	if _, err := h.users.UpsertHostProfile(r.Context(), claims.UserID, b.ID, body.DisplayName, body.Bio, body.City); err != nil {
		platform.WriteError(w, err)
		return
	}

	authResp, err := h.authClient.EnableHost(r.Context(), claims.UserID, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		platform.WriteError(w, platform.ErrInternal)
		return
	}

	platform.JSON(w, http.StatusOK, authResp)
}
