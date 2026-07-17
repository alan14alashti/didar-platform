// Package booking implements guest booking requests against experience
// time slots, and the host-side confirmation/cancellation workflow.
package booking

import "time"

type Booking struct {
	ID              string    `json:"id"`
	ExperienceID    string    `json:"experience_id"`
	ExperienceTitle string    `json:"experience_title,omitempty"`
	ExperienceSlug  string    `json:"experience_slug,omitempty"`
	CoverPhotoURL   *string   `json:"cover_photo_url,omitempty"`
	SlotID          string    `json:"slot_id"`
	SlotStartAt     time.Time `json:"slot_start_at,omitempty"`
	SlotEndAt       time.Time `json:"slot_end_at,omitempty"`
	GuestID         string    `json:"guest_id"`
	GuestName       *string   `json:"guest_name,omitempty"`
	GuestPhone      string    `json:"guest_phone,omitempty"`
	GuestsCount     int       `json:"guests_count"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	GuestNote       *string   `json:"guest_note"`
	HostNote        *string   `json:"host_note"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)
