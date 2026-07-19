package auth

import (
	"net/mail"
	"strings"
)

// NormalizeEmail lowercases and validates an email address.
// Returns "" when the input is not a valid email.
func NormalizeEmail(raw string) string {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return ""
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return ""
	}
	email := strings.TrimSpace(strings.ToLower(addr.Address))
	if !strings.Contains(email, "@") || strings.Contains(email, " ") {
		return ""
	}
	return email
}
