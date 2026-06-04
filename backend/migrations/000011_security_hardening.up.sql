-- Magic links: store token hash only; track failed login-code attempts.
-- Deploy cutover: expire outstanding links (old plaintext tokens are not carried forward).
ALTER TABLE magic_links
    ADD COLUMN IF NOT EXISTS token_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS failed_attempts INTEGER NOT NULL DEFAULT 0;

UPDATE magic_links
SET used_at = NOW()
WHERE used_at IS NULL;

ALTER TABLE magic_links DROP COLUMN IF EXISTS token;

DROP INDEX IF EXISTS idx_magic_links_token;
ALTER TABLE magic_links DROP CONSTRAINT IF EXISTS magic_links_token_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_magic_links_token_hash ON magic_links (token_hash)
    WHERE token_hash IS NOT NULL;
