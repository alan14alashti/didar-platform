-- Switch primary login identity to email; keep phone for a future SMS restore.
ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT;

UPDATE users SET email = lower(regexp_replace(phone, '^\+98', 'host') || '@didar.local')
WHERE email IS NULL AND phone IS NOT NULL;

ALTER TABLE users ALTER COLUMN phone DROP NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_unique ON users (phone) WHERE phone IS NOT NULL;

-- Drop the old NOT NULL UNIQUE constraint on phone if it still exists as a table constraint.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_phone_key;

-- Channel-agnostic OTP recipient (email or phone).
ALTER TABLE otp_codes RENAME COLUMN phone TO recipient;

-- Demo host emails (fixed UUIDs match seed).
UPDATE users SET email = 'arash@didar.local'    WHERE id = '11111111-1111-1111-1111-111111111101';
UPDATE users SET email = 'maryam@didar.local'   WHERE id = '11111111-1111-1111-1111-111111111102';
UPDATE users SET email = 'reza@didar.local'     WHERE id = '11111111-1111-1111-1111-111111111103';
UPDATE users SET email = 'negar@didar.local'    WHERE id = '11111111-1111-1111-1111-111111111104';
UPDATE users SET email = 'sina@didar.local'     WHERE id = '11111111-1111-1111-1111-111111111105';
UPDATE users SET email = 'fatemeh@didar.local'  WHERE id = '11111111-1111-1111-1111-111111111106';
