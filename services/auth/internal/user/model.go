package user

import "time"

type User struct {
	ID          string     `json:"id"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	FullName    *string    `json:"full_name"`
	AvatarURL   *string    `json:"avatar_url"`
	Locale      string     `json:"locale"`
	IsHost      bool       `json:"is_host"`
	IsAdmin     bool       `json:"is_admin"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
