CREATE TABLE developer_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Integration agent',
    key_prefix TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX developer_api_keys_user_active_idx
    ON developer_api_keys(user_id)
    WHERE revoked_at IS NULL;

CREATE INDEX developer_api_keys_key_hash_idx ON developer_api_keys(key_hash);
