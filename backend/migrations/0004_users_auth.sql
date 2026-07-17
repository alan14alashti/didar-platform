-- Users authenticate solely via phone number + OTP (no passwords).
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone           TEXT NOT NULL UNIQUE,
    full_name       TEXT,
    avatar_url      TEXT,
    locale          TEXT NOT NULL DEFAULT 'fa',
    is_host         BOOLEAN NOT NULL DEFAULT false,
    is_admin        BOOLEAN NOT NULL DEFAULT false,
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked')),
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE otp_codes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone           TEXT NOT NULL,
    purpose         TEXT NOT NULL DEFAULT 'login',
    code_hash       TEXT NOT NULL,
    attempt_count   INT NOT NULL DEFAULT 0,
    max_attempts    INT NOT NULL DEFAULT 5,
    expires_at      TIMESTAMPTZ NOT NULL,
    consumed_at     TIMESTAMPTZ,
    ip_address      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_otp_codes_phone_purpose ON otp_codes (phone, purpose, created_at DESC);

CREATE TABLE sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash  TEXT NOT NULL UNIQUE,
    user_agent          TEXT,
    ip_address          TEXT,
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);

CREATE TABLE host_profiles (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    brand_id        UUID REFERENCES brands(id),
    display_name    TEXT,
    bio             TEXT,
    city            TEXT,
    is_verified     BOOLEAN NOT NULL DEFAULT false,
    rating_avg      NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count    INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
