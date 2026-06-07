DROP INDEX IF EXISTS idx_users_avatar_reading_id;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_reading_id;
DROP TABLE IF EXISTS avatar_renders;
DROP TABLE IF EXISTS avatar_readings;
