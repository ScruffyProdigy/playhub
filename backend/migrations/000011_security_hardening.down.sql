DROP INDEX IF EXISTS idx_magic_links_token_hash;

ALTER TABLE magic_links
    ADD COLUMN IF NOT EXISTS token TEXT;

ALTER TABLE magic_links
    DROP COLUMN IF EXISTS failed_attempts,
    DROP COLUMN IF EXISTS token_hash;

CREATE INDEX IF NOT EXISTS idx_magic_links_token ON magic_links (token);
