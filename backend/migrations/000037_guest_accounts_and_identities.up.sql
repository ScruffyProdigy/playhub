-- Guest accounts, multi-email sign-in, and OIDC identity placeholders.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_guest BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS merged_into_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique
    ON users (email)
    WHERE email IS NOT NULL AND email <> '';

CREATE TABLE IF NOT EXISTS user_emails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_emails_email_unique UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_user_emails_user_id ON user_emails (user_id);

CREATE TABLE IF NOT EXISTS user_identities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(32) NOT NULL CHECK (provider IN ('google', 'discord', 'apple', 'facebook')),
    subject VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_identities_provider_subject_unique UNIQUE (provider, subject)
);

CREATE INDEX IF NOT EXISTS idx_user_identities_user_id ON user_identities (user_id);

INSERT INTO user_emails (user_id, email, is_primary, verified_at, created_at)
SELECT u.id, u.email, true, u.created_at, u.created_at
FROM users u
WHERE u.email IS NOT NULL
  AND u.email <> ''
ON CONFLICT (email) DO NOTHING;
