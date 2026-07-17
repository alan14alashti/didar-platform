package user

import "time"

type User struct {
	ID           string     `json:"id"`
	Phone        string     `json:"phone"`
	FullName     *string    `json:"full_name"`
	AvatarURL    *string    `json:"avatar_url"`
	Locale       string     `json:"locale"`
	IsHost       bool       `json:"is_host"`
	IsAdmin      bool       `json:"is_admin"`
	Status       string     `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type HostProfile struct {
	UserID      string    `json:"user_id"`
	DisplayName *string   `json:"display_name"`
	Bio         *string   `json:"bio"`
	City        *string   `json:"city"`
	IsVerified  bool      `json:"is_verified"`
	RatingAvg   float64   `json:"rating_avg"`
	RatingCount int       `json:"rating_count"`
	CreatedAt   time.Time `json:"created_at"`
}
