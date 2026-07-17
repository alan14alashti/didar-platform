CREATE TABLE experiences (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id        UUID NOT NULL REFERENCES brands(id),
    host_id         UUID NOT NULL REFERENCES users(id),
    category_id     UUID REFERENCES experience_categories(id),
    slug            TEXT NOT NULL UNIQUE,
    title           TEXT NOT NULL,
    summary         TEXT,
    description     TEXT,
    language        TEXT NOT NULL DEFAULT 'fa',
    price_amount    NUMERIC(12,0) NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'IRR',
    duration_minutes INT NOT NULL DEFAULT 120,
    max_guests      INT NOT NULL DEFAULT 6,
    min_guests      INT NOT NULL DEFAULT 1,
    city            TEXT NOT NULL,
    address         TEXT,
    location        GEOGRAPHY(Point, 4326) NOT NULL,
    status          TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    rating_avg      NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count    INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_experiences_location ON experiences USING GIST (location);
CREATE INDEX idx_experiences_status ON experiences (status);
CREATE INDEX idx_experiences_city ON experiences (city);
CREATE INDEX idx_experiences_host ON experiences (host_id);
CREATE INDEX idx_experiences_category ON experiences (category_id);

CREATE TABLE experience_photos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experience_id   UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    url             TEXT NOT NULL,
    alt             TEXT,
    sort_order      INT NOT NULL DEFAULT 0,
    is_cover        BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_experience_photos_experience ON experience_photos (experience_id);

CREATE TABLE experience_slots (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experience_id   UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    start_at        TIMESTAMPTZ NOT NULL,
    end_at          TIMESTAMPTZ NOT NULL,
    capacity        INT NOT NULL,
    booked_count    INT NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_experience_slots_experience ON experience_slots (experience_id, start_at);
