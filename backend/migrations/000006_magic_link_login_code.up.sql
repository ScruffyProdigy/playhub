ALTER TABLE magic_links
    ADD COLUMN code_hash VARCHAR(64);

CREATE INDEX idx_magic_links_email_active ON magic_links (email, created_at DESC)
    WHERE used_at IS NULL;
