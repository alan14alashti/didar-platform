// Package brand exposes the white-label configuration (name, logo,
// theme colors, locales, currency, domain) that drives the frontend's
// branding. Everything here is DB-backed and editable without a
// redeploy, which is what makes re-skinning the platform for another
// operator a config change rather than a code change.
package brand

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("brand not found")

type Brand struct {
	ID                string   `json:"id"`
	Slug              string   `json:"slug"`
	Name              string   `json:"name"`
	Tagline           *string  `json:"tagline"`
	Description       *string  `json:"description"`
	LogoURL           *string  `json:"logo_url"`
	FaviconURL        *string  `json:"favicon_url"`
	PrimaryColor      string   `json:"primary_color"`
	SecondaryColor    string   `json:"secondary_color"`
	AccentColor       string   `json:"accent_color"`
	DefaultLocale     string   `json:"default_locale"`
	SupportedLocales  []string `json:"supported_locales"`
	DefaultCurrency   string   `json:"default_currency"`
	ContactPhone      *string  `json:"contact_phone"`
	ContactEmail      *string  `json:"contact_email"`
	Domain            *string  `json:"domain"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const brandColumns = `id::text, slug, name, tagline, description, logo_url, favicon_url, primary_color, secondary_color, accent_color, default_locale, supported_locales, default_currency, contact_phone, contact_email, domain`

func scanBrand(row pgx.Row) (*Brand, error) {
	var b Brand
	if err := row.Scan(&b.ID, &b.Slug, &b.Name, &b.Tagline, &b.Description, &b.LogoURL, &b.FaviconURL, &b.PrimaryColor, &b.SecondaryColor, &b.AccentColor, &b.DefaultLocale, &b.SupportedLocales, &b.DefaultCurrency, &b.ContactPhone, &b.ContactEmail, &b.Domain); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*Brand, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+brandColumns+` FROM brands WHERE slug = $1 AND is_active = true`, slug)
	return scanBrand(row)
}

func (r *Repository) GetByDomain(ctx context.Context, domain string) (*Brand, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+brandColumns+` FROM brands WHERE domain = $1 AND is_active = true`, domain)
	return scanBrand(row)
}

// Resolve finds the active brand for an inbound request host, falling
// back to the platform default brand when no exact domain match exists
// (e.g. in local development where the host is "localhost:8080").
func (r *Repository) Resolve(ctx context.Context, host, defaultSlug string) (*Brand, error) {
	if host != "" {
		if b, err := r.GetByDomain(ctx, host); err == nil {
			return b, nil
		}
	}
	return r.GetBySlug(ctx, defaultSlug)
}
