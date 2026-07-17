package auth

import (
	"errors"
	"net/http"

	"github.com/didar-platform/backend/internal/brand"
	"github.com/didar-platform/backend/internal/otp"
	"github.com/didar-platform/backend/internal/platform"
	"github.com/didar-platform/backend/internal/user"
)

type Handler struct {
	users       *user.Repository
	otpService  *otp.Service
	issuer      *TokenIssuer
	sessions    *SessionRepository
	brands      *brand.Repository
	defaultSlug string
}

func NewHandler(users *user.Repository, otpService *otp.Service, issuer *TokenIssuer, sessions *SessionRepository, brands *brand.Repository, defaultBrandSlug string) *Handler {
	return &Handler{users: users, otpService: otpService, issuer: issuer, sessions: sessions, brands: brands, defaultSlug: defaultBrandSlug}
}

type requestOTPBody struct {
	Phone string `json:"phone"`
}

// RequestOTP sends a one-time code to the given Iranian mobile number
// via Kavenegar (or logs it in mock/dev mode).
func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var body requestOTPBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}

	phone := NormalizePhone(body.Phone)
	if phone == "" {
		platform.WriteError(w, platform.ErrValidation("Please enter a valid Iranian mobile number"))
		return
	}

	err := h.otpService.Request(r.Context(), phone, otp.PurposeLogin, r.RemoteAddr)
	switch {
	case errors.Is(err, otp.ErrTooManyRequests):
		platform.WriteError(w, platform.NewAPIError(http.StatusTooManyRequests, "otp_cooldown", "Please wait a bit before requesting another code"))
		return
	case err != nil:
		platform.WriteError(w, err)
		return
	}

	platform.JSON(w, http.StatusOK, map[string]any{"phone": phone, "sent": true})
}

type verifyOTPBody struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type authResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    string    `json:"expires_at"`
	User         *user.User `json:"user"`
	IsNewUser    bool      `json:"is_new_user"`
}

// VerifyOTP validates the code, creates the user on first login, and
// issues an access + refresh token pair.
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var body verifyOTPBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}

	phone := NormalizePhone(body.Phone)
	if phone == "" {
		platform.WriteError(w, platform.ErrValidation("Please enter a valid Iranian mobile number"))
		return
	}

	err := h.otpService.Verify(r.Context(), phone, otp.PurposeLogin, body.Code)
	switch {
	case errors.Is(err, otp.ErrIncorrectCode):
		platform.WriteError(w, platform.NewAPIError(http.StatusBadRequest, "otp_incorrect", "The code you entered is incorrect"))
		return
	case errors.Is(err, otp.ErrTooManyAttempts):
		platform.WriteError(w, platform.NewAPIError(http.StatusTooManyRequests, "otp_too_many_attempts", "Too many incorrect attempts, please request a new code"))
		return
	case errors.Is(err, otp.ErrExpiredOrUnknown):
		platform.WriteError(w, platform.NewAPIError(http.StatusBadRequest, "otp_expired", "This code has expired, please request a new one"))
		return
	case err != nil:
		platform.WriteError(w, err)
		return
	}

	isNewUser := false
	u, err := h.users.GetByPhone(r.Context(), phone)
	if errors.Is(err, user.ErrNotFound) {
		u, err = h.users.CreateWithPhone(r.Context(), phone)
		isNewUser = true
	} else if err == nil {
		_ = h.users.TouchLastLogin(r.Context(), u.ID)
	}
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	h.respondWithSession(w, r, u, isNewUser)
}

func (h *Handler) respondWithSession(w http.ResponseWriter, r *http.Request, u *user.User, isNewUser bool) {
	access, expiresAt, err := h.issuer.IssueAccessToken(u.ID, u.IsHost, u.IsAdmin)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	refresh, err := h.sessions.Create(r.Context(), u.ID, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.JSON(w, http.StatusOK, authResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:         u,
		IsNewUser:    isNewUser,
	})
}

type refreshBody struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body refreshBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	if body.RefreshToken == "" {
		platform.WriteError(w, platform.ErrValidation("refresh_token is required"))
		return
	}

	userID, newToken, err := h.sessions.Rotate(r.Context(), body.RefreshToken, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}

	u, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}

	access, expiresAt, err := h.issuer.IssueAccessToken(u.ID, u.IsHost, u.IsAdmin)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.JSON(w, http.StatusOK, authResponse{
		AccessToken:  access,
		RefreshToken: newToken,
		ExpiresAt:    expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:         u,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body refreshBody
	_ = platform.DecodeJSON(r, &body)
	if body.RefreshToken != "" {
		_ = h.sessions.Revoke(r.Context(), body.RefreshToken)
	}
	platform.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := FromContext(r.Context())
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	u, err := h.users.GetByID(r.Context(), claims.UserID)
	if err != nil {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	platform.JSON(w, http.StatusOK, u)
}

type updateMeBody struct {
	FullName *string `json:"full_name"`
	Locale   *string `json:"locale"`
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := FromContext(r.Context())
	if !ok {
		platform.WriteError(w, platform.ErrUnauthorized)
		return
	}
	var body updateMeBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}
	u, err := h.users.UpdateProfile(r.Context(), claims.UserID, body.FullName, body.Locale, nil)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.JSON(w, http.StatusOK, u)
}

type becomeHostBody struct {
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	City        *string `json:"city"`
}

// BecomeHost flips the caller into a host: from this point their JWTs
// (after re-login/refresh) carry is_host=true and the frontend shows the
// host dashboard alongside the regular guest experience.
func (h *Handler) BecomeHost(w http.ResponseWriter, r *http.Request) {
	claims, ok := FromContext(r.Context())
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

	u, err := h.users.BecomeHost(r.Context(), claims.UserID, b.ID, body.DisplayName, body.Bio, body.City)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	h.respondWithSession(w, r, u, false)
}
