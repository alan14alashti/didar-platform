-- Brands drive white-label theming/config so the same codebase can be
-- re-skinned for another operator without code changes.
CREATE TABLE brands (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug                TEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    tagline             TEXT,
    description         TEXT,
    logo_url            TEXT,
    favicon_url         TEXT,
    primary_color       TEXT NOT NULL DEFAULT '#FF385C',
    secondary_color     TEXT NOT NULL DEFAULT '#111827',
    accent_color        TEXT NOT NULL DEFAULT '#00A699',
    default_locale      TEXT NOT NULL DEFAULT 'fa',
    supported_locales   TEXT[] NOT NULL DEFAULT ARRAY['fa', 'en'],
    default_currency    TEXT NOT NULL DEFAULT 'IRR',
    contact_phone       TEXT,
    contact_email       TEXT,
    domain              TEXT UNIQUE,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
