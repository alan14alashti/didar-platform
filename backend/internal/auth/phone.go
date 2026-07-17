package auth

import (
	"regexp"
	"strings"
)

var iranMobileDigits = regexp.MustCompile(`^9\d{9}$`)

// NormalizePhone accepts common Iranian mobile input formats
// (09121234567, 9121234567, +989121234567, 00989121234567) and returns a
// canonical +98XXXXXXXXXX form, or "" if the input is not a valid
// Iranian mobile number.
func NormalizePhone(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")

	switch {
	case strings.HasPrefix(s, "+98"):
		s = strings.TrimPrefix(s, "+98")
	case strings.HasPrefix(s, "0098"):
		s = strings.TrimPrefix(s, "0098")
	case strings.HasPrefix(s, "98") && len(s) == 12:
		s = strings.TrimPrefix(s, "98")
	case strings.HasPrefix(s, "0"):
		s = strings.TrimPrefix(s, "0")
	}

	if !iranMobileDigits.MatchString(s) {
		return ""
	}
	return "+98" + s
}
