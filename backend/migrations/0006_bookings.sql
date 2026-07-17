CREATE TABLE bookings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experience_id   UUID NOT NULL REFERENCES experiences(id),
    slot_id         UUID NOT NULL REFERENCES experience_slots(id),
    guest_id        UUID NOT NULL REFERENCES users(id),
    guests_count    INT NOT NULL DEFAULT 1,
    total_amount    NUMERIC(12,0) NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'IRR',
    status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    guest_note      TEXT,
    host_note       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_bookings_guest ON bookings (guest_id);
CREATE INDEX idx_bookings_experience ON bookings (experience_id);
CREATE INDEX idx_bookings_slot ON bookings (slot_id);
