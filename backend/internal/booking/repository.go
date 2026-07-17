package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("booking not found")
	ErrSlotNotFound  = errors.New("slot not found")
	ErrSlotFull      = errors.New("slot has no remaining capacity")
	ErrSlotNotOpen   = errors.New("slot is not open for booking")
	ErrInvalidStatus = errors.New("invalid booking status transition")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create books guestsCount seats on a slot inside a single transaction,
// row-locking the slot to prevent overbooking under concurrent requests.
func (r *Repository) Create(ctx context.Context, guestID, slotID string, guestsCount int, guestNote *string) (*Booking, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var (
		experienceID string
		capacity     int
		bookedCount  int
		status       string
		priceAmount  float64
		currency     string
	)
	err = tx.QueryRow(ctx, `
		SELECT s.experience_id::text, s.capacity, s.booked_count, s.status, e.price_amount, e.currency
		FROM experience_slots s
		JOIN experiences e ON e.id = s.experience_id
		WHERE s.id = $1
		FOR UPDATE OF s
	`, slotID).Scan(&experienceID, &capacity, &bookedCount, &status, &priceAmount, &currency)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSlotNotFound
		}
		return nil, fmt.Errorf("load slot: %w", err)
	}

	if status != "open" {
		return nil, ErrSlotNotOpen
	}
	if bookedCount+guestsCount > capacity {
		return nil, ErrSlotFull
	}

	if _, err := tx.Exec(ctx, `UPDATE experience_slots SET booked_count = booked_count + $1 WHERE id = $2`, guestsCount, slotID); err != nil {
		return nil, fmt.Errorf("update slot capacity: %w", err)
	}

	total := priceAmount * float64(guestsCount)

	var b Booking
	err = tx.QueryRow(ctx, `
		INSERT INTO bookings (experience_id, slot_id, guest_id, guests_count, total_amount, currency, guest_note)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, experience_id::text, slot_id::text, guest_id::text, guests_count, total_amount, currency, status, guest_note, host_note, created_at, updated_at
	`, experienceID, slotID, guestID, guestsCount, total, currency, guestNote).Scan(
		&b.ID, &b.ExperienceID, &b.SlotID, &b.GuestID, &b.GuestsCount, &b.TotalAmount, &b.Currency, &b.Status, &b.GuestNote, &b.HostNote, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &b, nil
}

const guestBookingSelect = `
	SELECT
		b.id::text, b.experience_id::text, e.title, e.slug,
		(SELECT p.url FROM experience_photos p WHERE p.experience_id = e.id ORDER BY p.is_cover DESC, p.sort_order LIMIT 1),
		b.slot_id::text, s.start_at, s.end_at,
		b.guest_id::text, b.guests_count, b.total_amount, b.currency, b.status, b.guest_note, b.host_note, b.created_at, b.updated_at
	FROM bookings b
	JOIN experiences e ON e.id = b.experience_id
	JOIN experience_slots s ON s.id = b.slot_id
`

func scanGuestBooking(row pgx.Row) (Booking, error) {
	var b Booking
	err := row.Scan(&b.ID, &b.ExperienceID, &b.ExperienceTitle, &b.ExperienceSlug, &b.CoverPhotoURL,
		&b.SlotID, &b.SlotStartAt, &b.SlotEndAt,
		&b.GuestID, &b.GuestsCount, &b.TotalAmount, &b.Currency, &b.Status, &b.GuestNote, &b.HostNote, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *Repository) ListByGuest(ctx context.Context, guestID string) ([]Booking, error) {
	rows, err := r.pool.Query(ctx, guestBookingSelect+` WHERE b.guest_id = $1 ORDER BY s.start_at DESC`, guestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Booking, 0)
	for rows.Next() {
		b, err := scanGuestBooking(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}

func (r *Repository) GetByIDForGuest(ctx context.Context, id, guestID string) (*Booking, error) {
	row := r.pool.QueryRow(ctx, guestBookingSelect+` WHERE b.id = $1 AND b.guest_id = $2`, id, guestID)
	b, err := scanGuestBooking(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

const hostBookingSelect = `
	SELECT
		b.id::text, b.experience_id::text, e.title, e.slug,
		(SELECT p.url FROM experience_photos p WHERE p.experience_id = e.id ORDER BY p.is_cover DESC, p.sort_order LIMIT 1),
		b.slot_id::text, s.start_at, s.end_at,
		b.guest_id::text, u.full_name, u.phone, b.guests_count, b.total_amount, b.currency, b.status, b.guest_note, b.host_note, b.created_at, b.updated_at
	FROM bookings b
	JOIN experiences e ON e.id = b.experience_id
	JOIN experience_slots s ON s.id = b.slot_id
	JOIN users u ON u.id = b.guest_id
`

func scanHostBooking(row pgx.Row) (Booking, error) {
	var b Booking
	err := row.Scan(&b.ID, &b.ExperienceID, &b.ExperienceTitle, &b.ExperienceSlug, &b.CoverPhotoURL,
		&b.SlotID, &b.SlotStartAt, &b.SlotEndAt,
		&b.GuestID, &b.GuestName, &b.GuestPhone, &b.GuestsCount, &b.TotalAmount, &b.Currency, &b.Status, &b.GuestNote, &b.HostNote, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *Repository) ListByHost(ctx context.Context, hostID string) ([]Booking, error) {
	rows, err := r.pool.Query(ctx, hostBookingSelect+` WHERE e.host_id = $1 ORDER BY s.start_at DESC`, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Booking, 0)
	for rows.Next() {
		b, err := scanHostBooking(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}

// UpdateStatusByHost transitions a booking's status, releasing slot
// capacity when a booking is cancelled.
func (r *Repository) UpdateStatusByHost(ctx context.Context, bookingID, hostID, newStatus string, hostNote *string) (*Booking, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var (
		slotID       string
		guestsCount  int
		currentStatus string
	)
	err = tx.QueryRow(ctx, `
		SELECT b.slot_id::text, b.guests_count, b.status
		FROM bookings b
		JOIN experiences e ON e.id = b.experience_id
		WHERE b.id = $1 AND e.host_id = $2
		FOR UPDATE OF b
	`, bookingID, hostID).Scan(&slotID, &guestsCount, &currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if currentStatus == StatusCancelled || currentStatus == StatusCompleted {
		return nil, ErrInvalidStatus
	}

	if _, err := tx.Exec(ctx, `UPDATE bookings SET status = $2, host_note = COALESCE($3, host_note), updated_at = now() WHERE id = $1`, bookingID, newStatus, hostNote); err != nil {
		return nil, err
	}

	if newStatus == StatusCancelled {
		if _, err := tx.Exec(ctx, `UPDATE experience_slots SET booked_count = GREATEST(booked_count - $1, 0) WHERE id = $2`, guestsCount, slotID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, hostBookingSelect+` WHERE b.id = $1`, bookingID)
	b, err := scanHostBooking(row)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// CancelByGuest lets a guest cancel their own pending/confirmed booking.
func (r *Repository) CancelByGuest(ctx context.Context, bookingID, guestID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var slotID string
	var guestsCount int
	var status string
	err = tx.QueryRow(ctx, `
		SELECT slot_id::text, guests_count, status FROM bookings WHERE id = $1 AND guest_id = $2 FOR UPDATE
	`, bookingID, guestID).Scan(&slotID, &guestsCount, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if status == StatusCancelled || status == StatusCompleted {
		return ErrInvalidStatus
	}

	if _, err := tx.Exec(ctx, `UPDATE bookings SET status = $2, updated_at = now() WHERE id = $1`, bookingID, StatusCancelled); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE experience_slots SET booked_count = GREATEST(booked_count - $1, 0) WHERE id = $2`, guestsCount, slotID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
