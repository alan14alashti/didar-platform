// Package platform holds small cross-cutting HTTP helpers shared by all
// feature packages (JSON encoding/decoding, error envelopes, context keys).
package platform

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// APIError is a machine-readable error the frontend can key off of
// (e.g. "otp_expired", "phone_invalid") plus a human-readable message.
type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

func NewAPIError(status int, code, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}

var (
	ErrValidation   = func(msg string) *APIError { return NewAPIError(http.StatusBadRequest, "validation_error", msg) }
	ErrUnauthorized = NewAPIError(http.StatusUnauthorized, "unauthorized", "Authentication required")
	ErrForbidden    = NewAPIError(http.StatusForbidden, "forbidden", "You are not allowed to do this")
	ErrNotFound     = NewAPIError(http.StatusNotFound, "not_found", "Resource not found")
	ErrInternal     = NewAPIError(http.StatusInternalServerError, "internal_error", "Something went wrong")
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode json response", "error", err)
	}
}

func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return NewAPIError(http.StatusBadRequest, "invalid_body", "Request body is not valid JSON")
	}
	return nil
}

// WriteError inspects err and writes the appropriate JSON error envelope.
func WriteError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		JSON(w, apiErr.Status, map[string]any{"error": apiErr})
		return
	}
	slog.Error("unhandled error", "error", err)
	JSON(w, ErrInternal.Status, map[string]any{"error": ErrInternal})
}
