package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	serviceKey string
	httpClient *http.Client
}

func New(baseURL, serviceKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		serviceKey: serviceKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type User struct {
	ID       string  `json:"id"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	FullName *string `json:"full_name"`
	IsHost   bool    `json:"is_host"`
	IsAdmin  bool    `json:"is_admin"`
	Locale   string  `json:"locale"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	User         *User  `json:"user"`
	IsNewUser    bool   `json:"is_new_user"`
}

func (c *Client) EnableHost(ctx context.Context, userID, userAgent, remoteAddr string) (*AuthResponse, error) {
	url := c.baseURL + "/internal/users/" + userID + "/enable-host"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("build enable-host request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", c.serviceKey)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	if remoteAddr != "" {
		req.Header.Set("X-Forwarded-For", remoteAddr)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call auth service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned %d: %s", resp.StatusCode, string(body))
	}

	var out AuthResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode auth response: %w", err)
	}
	return &out, nil
}
