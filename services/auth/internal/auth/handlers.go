package auth

import (
	"errors"
	"net/http"

	"github.com/didar-platform/pkg/authkit"
	"github.com/didar-platform/pkg/platform"
	"github.com/didar-platform/services/auth/internal/otp"
	"github.com/didar-platform/services/auth/internal/user"
)

type Handler struct {
	users      *user.Repository
	otpService *otp.Service
	issuer     *authkit.TokenIssuer
	sessions   *SessionRepository
}

func NewHandler(users *user.Repository, otpService *otp.Service, issuer *authkit.TokenIssuer, sessions *SessionRepository) *Handler {
	return &Handler{users: users, otpService: otpService, issuer: issuer, sessions: sessions}
}

type requestOTPBody struct {
	Email string `json:"email"`
}

func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var body requestOTPBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}

	email := NormalizeEmail(body.Email)
	if email == "" {
		platform.WriteError(w, platform.ErrValidation("Please enter a valid email address"))
		return
	}

	err := h.otpService.Request(r.Context(), email, otp.PurposeLogin, r.RemoteAddr)
	switch {
	case errors.Is(err, otp.ErrTooManyRequests):
		platform.WriteError(w, platform.NewAPIError(http.StatusTooManyRequests, "otp_cooldown", "Please wait a bit before requesting another code"))
		return
	case err != nil:
		platform.WriteError(w, err)
		return
	}

	platform.JSON(w, http.StatusOK, map[string]any{"email": email, "sent": true})
}

type verifyOTPBody struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    string     `json:"expires_at"`
	User         *user.User `json:"user"`
	IsNewUser    bool       `json:"is_new_user"`
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var body verifyOTPBody
	if err := platform.DecodeJSON(r, &body); err != nil {
		platform.WriteError(w, err)
		return
	}

	email := NormalizeEmail(body.Email)
	if email == "" {
		platform.WriteError(w, platform.ErrValidation("Please enter a valid email address"))
		return
	}

	err := h.otpService.Verify(r.Context(), email, otp.PurposeLogin, body.Code)
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
	u, err := h.users.GetByEmail(r.Context(), email)
	if errors.Is(err, user.ErrNotFound) {
		u, err = h.users.CreateWithEmail(r.Context(), email)
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

	platform.JSON(w, http.StatusOK, AuthResponse{
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

	platform.JSON(w, http.StatusOK, AuthResponse{
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
	claims, ok := authkit.FromContext(r.Context())
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
	claims, ok := authkit.FromContext(r.Context())
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

func (h *Handler) EnableHost(w http.ResponseWriter, r *http.Request, userID string) {
	u, err := h.users.EnableHost(r.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			platform.WriteError(w, platform.ErrNotFound)
			return
		}
		platform.WriteError(w, err)
		return
	}
	h.respondWithSession(w, r, u, false)
}
