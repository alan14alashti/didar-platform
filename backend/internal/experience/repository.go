package experience

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("experience not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT id::text, slug, title_fa, title_en, icon, sort_order FROM experience_categories ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Slug, &c.TitleFa, &c.TitleEn, &c.Icon, &c.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

const baseExperienceSelect = `
	SELECT
		e.id::text, e.brand_id::text, e.host_id::text, e.category_id::text, c.slug AS category_slug, e.slug, e.title, e.summary, e.description,
		e.language, e.price_amount, e.currency, e.duration_minutes, e.max_guests, e.min_guests, e.city, e.address,
		ST_Y(e.location::geometry) AS latitude, ST_X(e.location::geometry) AS longitude,
		e.status, e.rating_avg, e.rating_count, e.created_at, e.updated_at,
		(SELECT p.url FROM experience_photos p WHERE p.experience_id = e.id ORDER BY p.is_cover DESC, p.sort_order LIMIT 1) AS cover_photo_url
`

func scanExperience(row pgx.Row, withDistance bool, distancePtr *float64) (Experience, error) {
	var e Experience
	dest := []any{
		&e.ID, &e.BrandID, &e.HostID, &e.CategoryID, &e.CategorySlug, &e.Slug, &e.Title, &e.Summary, &e.Description,
		&e.Language, &e.PriceAmount, &e.Currency, &e.DurationMinutes, &e.MaxGuests, &e.MinGuests, &e.City, &e.Address,
		&e.Latitude, &e.Longitude,
		&e.Status, &e.RatingAvg, &e.RatingCount, &e.CreatedAt, &e.UpdatedAt, &e.CoverPhotoURL,
	}
	if withDistance {
		dest = append(dest, distancePtr)
	}
	if err := row.Scan(dest...); err != nil {
		return e, err
	}
	if withDistance {
		e.DistanceMeters = distancePtr
	}
	return e, nil
}

// Search runs a filtered, paginated experience listing. When Lat/Lng are
// provided, results are ordered by proximity using PostGIS' geography
// distance operator, which is index-accelerated via the GIST index on
// experiences.location.
func (r *Repository) Search(ctx context.Context, p SearchParams) (*SearchResult, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 50 {
		p.PageSize = 12
	}

	var (
		conditions []string
		args       []any
	)
	conditions = append(conditions, "e.status = 'published'")

	addArg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	if p.Query != "" {
		ph := addArg("%" + strings.ToLower(p.Query) + "%")
		conditions = append(conditions, fmt.Sprintf("(lower(e.title) LIKE %s OR lower(e.summary) LIKE %s OR lower(e.city) LIKE %s)", ph, ph, ph))
	}
	if p.CategorySlug != "" {
		ph := addArg(p.CategorySlug)
		conditions = append(conditions, fmt.Sprintf("c.slug = %s", ph))
	}
	if p.City != "" {
		ph := addArg(p.City)
		conditions = append(conditions, fmt.Sprintf("e.city = %s", ph))
	}

	var geoSelect, geoOrder string
	if p.Lat != nil && p.Lng != nil {
		lngPh := addArg(*p.Lng)
		latPh := addArg(*p.Lat)
		point := fmt.Sprintf("ST_SetSRID(ST_MakePoint(%s, %s), 4326)::geography", lngPh, latPh)
		if p.RadiusKM > 0 {
			radiusPh := addArg(p.RadiusKM * 1000)
			conditions = append(conditions, fmt.Sprintf("ST_DWithin(e.location, %s, %s)", point, radiusPh))
		}
		geoSelect = fmt.Sprintf(", ST_Distance(e.location, %s) AS distance_meters", point)
		geoOrder = "distance_meters ASC,"
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT count(*) FROM experiences e
		LEFT JOIN experience_categories c ON c.id = e.category_id
		WHERE %s
	`, where)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count experiences: %w", err)
	}

	limitPh := addArg(p.PageSize)
	offsetPh := addArg((p.Page - 1) * p.PageSize)

	query := fmt.Sprintf(`
		%s
		%s
		FROM experiences e
		LEFT JOIN experience_categories c ON c.id = e.category_id
		WHERE %s
		ORDER BY %s e.created_at DESC
		LIMIT %s OFFSET %s
	`, baseExperienceSelect, geoSelect, where, geoOrder, limitPh, offsetPh)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("search experiences: %w", err)
	}
	defer rows.Close()

	withDistance := geoSelect != ""
	items := make([]Experience, 0)
	for rows.Next() {
		var dist float64
		e, err := scanExperience(rows, withDistance, &dist)
		if err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &SearchResult{Items: items, Page: p.Page, PageSize: p.PageSize, TotalCount: total}, nil
}

func (r *Repository) GetBySlugPublished(ctx context.Context, slug string) (*Experience, error) {
	query := baseExperienceSelect + `
		FROM experiences e
		LEFT JOIN experience_categories c ON c.id = e.category_id
		WHERE e.slug = $1 AND e.status = 'published'
	`
	row := r.pool.QueryRow(ctx, query, slug)
	e, err := scanExperience(row, false, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := r.attachPhotosAndSlots(ctx, &e); err != nil {
		return nil, err
	}
	if err := r.attachHost(ctx, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *Repository) attachPhotosAndSlots(ctx context.Context, e *Experience) error {
	photoRows, err := r.pool.Query(ctx, `SELECT id::text, url, alt, sort_order, is_cover FROM experience_photos WHERE experience_id = $1 ORDER BY is_cover DESC, sort_order`, e.ID)
	if err != nil {
		return err
	}
	defer photoRows.Close()
	for photoRows.Next() {
		var p Photo
		if err := photoRows.Scan(&p.ID, &p.URL, &p.Alt, &p.SortOrder, &p.IsCover); err != nil {
			return err
		}
		e.Photos = append(e.Photos, p)
	}
	if err := photoRows.Err(); err != nil {
		return err
	}

	slotRows, err := r.pool.Query(ctx, `
		SELECT id::text, experience_id::text, start_at, end_at, capacity, booked_count, status
		FROM experience_slots
		WHERE experience_id = $1 AND start_at > now() AND status = 'open'
		ORDER BY start_at
	`, e.ID)
	if err != nil {
		return err
	}
	defer slotRows.Close()
	for slotRows.Next() {
		var s Slot
		if err := slotRows.Scan(&s.ID, &s.ExperienceID, &s.StartAt, &s.EndAt, &s.Capacity, &s.BookedCount, &s.Status); err != nil {
			return err
		}
		e.Slots = append(e.Slots, s)
	}
	return slotRows.Err()
}

func (r *Repository) attachHost(ctx context.Context, e *Experience) error {
	row := r.pool.QueryRow(ctx, `
		SELECT u.id::text, u.full_name, u.avatar_url, hp.display_name, hp.bio, hp.is_verified, hp.rating_avg, hp.rating_count
		FROM users u
		LEFT JOIN host_profiles hp ON hp.user_id = u.id
		WHERE u.id = $1
	`, e.HostID)
	var h HostSummary
	if err := row.Scan(&h.ID, &h.FullName, &h.AvatarURL, &h.DisplayName, &h.Bio, &h.IsVerified, &h.RatingAvg, &h.RatingCount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	e.Host = &h
	return nil
}

// --- Host-owned CRUD ---

type UpsertExperienceInput struct {
	CategoryID      *string
	Title           string
	Summary         *string
	Description     *string
	PriceAmount     float64
	DurationMinutes int
	MaxGuests       int
	MinGuests       int
	City            string
	Address         *string
	Latitude        float64
	Longitude       float64
	Status          string
}

func (r *Repository) ListByHost(ctx context.Context, hostID string) ([]Experience, error) {
	query := baseExperienceSelect + `
		FROM experiences e
		LEFT JOIN experience_categories c ON c.id = e.category_id
		WHERE e.host_id = $1
		ORDER BY e.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Experience, 0)
	for rows.Next() {
		e, err := scanExperience(rows, false, nil)
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *Repository) GetByIDForHost(ctx context.Context, id, hostID string) (*Experience, error) {
	query := baseExperienceSelect + `
		FROM experiences e
		LEFT JOIN experience_categories c ON c.id = e.category_id
		WHERE e.id = $1 AND e.host_id = $2
	`
	row := r.pool.QueryRow(ctx, query, id, hostID)
	e, err := scanExperience(row, false, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := r.attachPhotosAndSlots(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func slugify(title string) string {
	// Keep it simple: lowercased, spaces to dashes, plus a short random
	// suffix appended by the caller to guarantee uniqueness for
	// non-Latin (e.g. Persian) titles which have no ASCII slug.
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "experience"
	}
	return out
}

func (r *Repository) Create(ctx context.Context, brandID, hostID string, in UpsertExperienceInput) (*Experience, error) {
	baseSlug := slugify(in.Title)
	slug := baseSlug

	var id string
	for attempt := 0; attempt < 5; attempt++ {
		err := r.pool.QueryRow(ctx, `
			INSERT INTO experiences (
				brand_id, host_id, category_id, slug, title, summary, description,
				price_amount, duration_minutes, max_guests, min_guests, city, address, location, status
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
				ST_SetSRID(ST_MakePoint($14, $15), 4326)::geography, $16
			) RETURNING id::text
		`, brandID, hostID, in.CategoryID, slug, in.Title, in.Summary, in.Description,
			in.PriceAmount, in.DurationMinutes, in.MaxGuests, in.MinGuests, in.City, in.Address,
			in.Longitude, in.Latitude, in.Status,
		).Scan(&id)
		if err == nil {
			return r.GetByIDForHost(ctx, id, hostID)
		}
		if isUniqueViolation(err) {
			slug = fmt.Sprintf("%s-%d", baseSlug, attempt+2)
			continue
		}
		return nil, fmt.Errorf("create experience: %w", err)
	}
	return nil, errors.New("could not generate a unique slug")
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "SQLSTATE 23505")
}

func (r *Repository) Update(ctx context.Context, id, hostID string, in UpsertExperienceInput) (*Experience, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE experiences SET
			category_id = $3, title = $4, summary = $5, description = $6, price_amount = $7,
			duration_minutes = $8, max_guests = $9, min_guests = $10, city = $11, address = $12,
			location = ST_SetSRID(ST_MakePoint($13, $14), 4326)::geography, status = $15, updated_at = now()
		WHERE id = $1 AND host_id = $2
	`, id, hostID, in.CategoryID, in.Title, in.Summary, in.Description, in.PriceAmount,
		in.DurationMinutes, in.MaxGuests, in.MinGuests, in.City, in.Address,
		in.Longitude, in.Latitude, in.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("update experience: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.GetByIDForHost(ctx, id, hostID)
}

func (r *Repository) Archive(ctx context.Context, id, hostID string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE experiences SET status = 'archived', updated_at = now() WHERE id = $1 AND host_id = $2`, id, hostID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) AddPhoto(ctx context.Context, experienceID, hostID, url string, alt *string, isCover bool) (*Photo, error) {
	var owns bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM experiences WHERE id = $1 AND host_id = $2)`, experienceID, hostID).Scan(&owns); err != nil {
		return nil, err
	}
	if !owns {
		return nil, ErrNotFound
	}

	if isCover {
		if _, err := r.pool.Exec(ctx, `UPDATE experience_photos SET is_cover = false WHERE experience_id = $1`, experienceID); err != nil {
			return nil, err
		}
	}

	var nextSort int
	if err := r.pool.QueryRow(ctx, `SELECT COALESCE(MAX(sort_order) + 1, 0) FROM experience_photos WHERE experience_id = $1`, experienceID).Scan(&nextSort); err != nil {
		return nil, err
	}

	var p Photo
	err := r.pool.QueryRow(ctx, `
		INSERT INTO experience_photos (experience_id, url, alt, sort_order, is_cover)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, url, alt, sort_order, is_cover
	`, experienceID, url, alt, nextSort, isCover).Scan(&p.ID, &p.URL, &p.Alt, &p.SortOrder, &p.IsCover)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) DeletePhoto(ctx context.Context, experienceID, hostID, photoID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM experience_photos
		WHERE id = $1 AND experience_id = $2
		AND EXISTS (SELECT 1 FROM experiences WHERE id = $2 AND host_id = $3)
	`, photoID, experienceID, hostID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) AddSlot(ctx context.Context, experienceID, hostID string, s Slot) (*Slot, error) {
	var owns bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM experiences WHERE id = $1 AND host_id = $2)`, experienceID, hostID).Scan(&owns); err != nil {
		return nil, err
	}
	if !owns {
		return nil, ErrNotFound
	}

	var out Slot
	err := r.pool.QueryRow(ctx, `
		INSERT INTO experience_slots (experience_id, start_at, end_at, capacity)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, experience_id::text, start_at, end_at, capacity, booked_count, status
	`, experienceID, s.StartAt, s.EndAt, s.Capacity).Scan(&out.ID, &out.ExperienceID, &out.StartAt, &out.EndAt, &out.Capacity, &out.BookedCount, &out.Status)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) DeleteSlot(ctx context.Context, experienceID, hostID, slotID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM experience_slots
		WHERE id = $1 AND experience_id = $2 AND booked_count = 0
		AND EXISTS (SELECT 1 FROM experiences WHERE id = $2 AND host_id = $3)
	`, slotID, experienceID, hostID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
