DROP INDEX IF EXISTS idx_magic_links_email_active;

ALTER TABLE magic_links
    DROP COLUMN IF EXISTS code_hash;
