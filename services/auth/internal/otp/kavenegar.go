package otp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kavenegar/kavenegar-go"
)

// KavenegarSender sends OTP codes via Kavenegar Verify Lookup.
// Kept for restoring phone SMS auth later.
type KavenegarSender struct {
	api          *kavenegar.Kavenegar
	templateName string
}

func NewKavenegarSender(apiKey, templateName, _ string) *KavenegarSender {
	return &KavenegarSender{
		api:          kavenegar.New(apiKey),
		templateName: templateName,
	}
}

func (k *KavenegarSender) SendOTP(_ context.Context, phone, code string) error {
	receptor := normalizeForKavenegar(phone)
	params := &kavenegar.VerifyLookupParam{}

	res, err := k.api.Verify.Lookup(receptor, k.templateName, code, params)
	if err != nil {
		switch e := err.(type) {
		case *kavenegar.APIError:
			return fmt.Errorf("kavenegar api error: %s", e.Error())
		case *kavenegar.HTTPError:
			return fmt.Errorf("kavenegar http error: %s", e.Error())
		default:
			return fmt.Errorf("kavenegar error: %w", err)
		}
	}

	slog.Info("kavenegar otp sent",
		"phone", receptor,
		"message_id", res.MessageID,
		"status", res.Status,
	)
	return nil
}

func normalizeForKavenegar(phone string) string {
	p := strings.TrimSpace(phone)
	switch {
	case strings.HasPrefix(p, "+98"):
		return "0" + strings.TrimPrefix(p, "+98")
	case strings.HasPrefix(p, "98") && len(p) > 10:
		return "0" + strings.TrimPrefix(p, "98")
	default:
		return p
	}
}
