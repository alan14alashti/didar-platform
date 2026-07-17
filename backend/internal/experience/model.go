// Package experience covers the core marketplace listing: categories,
// experiences (with a PostGIS point for geo search), their photo
// galleries and bookable time slots.
package experience

import "time"

type Category struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	TitleFa   string `json:"title_fa"`
	TitleEn   string `json:"title_en"`
	Icon      *string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

type Photo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Alt       *string `json:"alt"`
	SortOrder int    `json:"sort_order"`
	IsCover   bool   `json:"is_cover"`
}

type Slot struct {
	ID          string    `json:"id"`
	ExperienceID string   `json:"experience_id"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Capacity    int       `json:"capacity"`
	BookedCount int       `json:"booked_count"`
	Status      string    `json:"status"`
}

func (s Slot) SeatsLeft() int {
	left := s.Capacity - s.BookedCount
	if left < 0 {
		return 0
	}
	return left
}

type HostSummary struct {
	ID          string  `json:"id"`
	FullName    *string `json:"full_name"`
	AvatarURL   *string `json:"avatar_url"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	IsVerified  bool    `json:"is_verified"`
	RatingAvg   float64 `json:"rating_avg"`
	RatingCount int     `json:"rating_count"`
}

type Experience struct {
	ID              string    `json:"id"`
	BrandID         string    `json:"brand_id"`
	HostID          string    `json:"host_id"`
	CategoryID      *string   `json:"category_id"`
	CategorySlug    *string   `json:"category_slug,omitempty"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Summary         *string   `json:"summary"`
	Description     *string   `json:"description"`
	Language        string    `json:"language"`
	PriceAmount     float64   `json:"price_amount"`
	Currency        string    `json:"currency"`
	DurationMinutes int       `json:"duration_minutes"`
	MaxGuests       int       `json:"max_guests"`
	MinGuests       int       `json:"min_guests"`
	City            string    `json:"city"`
	Address         *string   `json:"address"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Status          string    `json:"status"`
	RatingAvg       float64   `json:"rating_avg"`
	RatingCount     int       `json:"rating_count"`
	DistanceMeters  *float64  `json:"distance_meters,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	CoverPhotoURL *string     `json:"cover_photo_url,omitempty"`
	Photos        []Photo     `json:"photos,omitempty"`
	Slots         []Slot      `json:"slots,omitempty"`
	Host          *HostSummary `json:"host,omitempty"`
}

type SearchParams struct {
	Query      string
	CategorySlug string
	City       string
	Lat        *float64
	Lng        *float64
	RadiusKM   float64
	Page       int
	PageSize   int
}

type SearchResult struct {
	Items      []Experience `json:"items"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalCount int          `json:"total_count"`
}
