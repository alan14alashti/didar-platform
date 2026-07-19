-- Move auth tables to the auth microservice database.
-- User IDs remain as opaque UUID references in this database.

ALTER TABLE IF EXISTS experiences DROP CONSTRAINT IF EXISTS experiences_host_id_fkey;
ALTER TABLE IF EXISTS bookings DROP CONSTRAINT IF EXISTS bookings_guest_id_fkey;
ALTER TABLE IF EXISTS host_profiles DROP CONSTRAINT IF EXISTS host_profiles_user_id_fkey;

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS otp_codes;
DROP TABLE IF EXISTS users;
