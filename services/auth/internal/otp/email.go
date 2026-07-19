package otp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

// EmailSender delivers OTP codes over SMTP.
type EmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailSender(host string, port int, username, password, from string) *EmailSender {
	if port <= 0 {
		port = 587
	}
	return &EmailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (e *EmailSender) SendOTP(_ context.Context, recipient, code string) error {
	from := e.From
	if from == "" {
		from = e.Username
	}
	if from == "" {
		return fmt.Errorf("smtp from address is not configured")
	}

	subject := "Didar verification code"
	body := fmt.Sprintf("Your Didar verification code is: %s\n\nThis code expires shortly. If you did not request it, ignore this email.\n", code)
	msg := strings.Join([]string{
		"From: " + from,
		"To: " + recipient,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := net.JoinHostPort(e.Host, strconv.Itoa(e.Port))
	var auth smtp.Auth
	if e.Username != "" {
		auth = smtp.PlainAuth("", e.Username, e.Password, e.Host)
	}

	if err := smtp.SendMail(addr, auth, from, []string{recipient}, []byte(msg)); err != nil {
		return fmt.Errorf("send email otp: %w", err)
	}

	slog.Info("email otp sent", "email", recipient)
	return nil
}

// MockSender logs the OTP instead of sending it.
type MockSender struct{}

func (MockSender) SendOTP(_ context.Context, recipient, code string) error {
	slog.Warn("OTP mock mode: code not actually sent", "recipient", recipient, "code", code)
	return nil
}
