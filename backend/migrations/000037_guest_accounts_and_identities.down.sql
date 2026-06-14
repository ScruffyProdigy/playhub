DROP TABLE IF EXISTS user_identities;
DROP TABLE IF EXISTS user_emails;

ALTER TABLE users DROP COLUMN IF EXISTS merged_into_user_id;
ALTER TABLE users DROP COLUMN IF EXISTS is_guest;

DROP INDEX IF EXISTS idx_users_email_unique;
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
