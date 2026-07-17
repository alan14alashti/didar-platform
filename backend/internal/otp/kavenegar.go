// Package otp implements one-time-password generation, storage and
// delivery over SMS via Kavenegar (https://kavenegar.com), the most
// widely used SMS/OTP provider for Iranian phone numbers.
package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SMSSender delivers an OTP code to a phone number. Swappable so tests
// and local development can use a mock instead of calling out to Kavenegar.
type SMSSender interface {
	SendOTP(ctx context.Context, phone, code string) error
}

// KavenegarSender calls Kavenegar's "Verify Lookup" endpoint, the
// recommended way to send OTP codes in Iran (pre-approved template,
// higher deliverability, no need to author full SMS text per-request).
type KavenegarSender struct {
	APIKey       string
	TemplateName string
	SenderLine   string
	HTTPClient   *http.Client
}

func NewKavenegarSender(apiKey, templateName, senderLine string) *KavenegarSender {
	return &KavenegarSender{
		APIKey:       apiKey,
		TemplateName: templateName,
		SenderLine:   senderLine,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

type kavenegarResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
}

func (k *KavenegarSender) SendOTP(ctx context.Context, phone, code string) error {
	endpoint := fmt.Sprintf("https://api.kavenegar.com/v1/%s/verify/lookup.json", url.PathEscape(k.APIKey))

	form := url.Values{}
	form.Set("receptor", normalizeForKavenegar(phone))
	form.Set("token", code)
	form.Set("template", k.TemplateName)
	if k.SenderLine != "" {
		form.Set("sender", k.SenderLine)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build kavenegar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("call kavenegar: %w", err)
	}
	defer resp.Body.Close()

	var parsed kavenegarResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("decode kavenegar response: %w", err)
	}

	if parsed.Return.Status != 200 {
		return fmt.Errorf("kavenegar error %d: %s", parsed.Return.Status, parsed.Return.Message)
	}

	return nil
}

// normalizeForKavenegar converts +98XXXXXXXXXX to the 0XXXXXXXXXX format
// Kavenegar/Iranian carriers expect.
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

// MockSender logs the OTP instead of sending it. Used automatically in
// development when no Kavenegar API key is configured.
type MockSender struct{}

func (MockSender) SendOTP(_ context.Context, phone, code string) error {
	slog.Warn("OTP mock mode: SMS not actually sent", "phone", phone, "code", code)
	return nil
}
