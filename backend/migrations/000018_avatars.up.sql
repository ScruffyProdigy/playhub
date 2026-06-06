ALTER TABLE users
    ADD COLUMN avatar_key VARCHAR(50),
    ADD COLUMN avatar_source VARCHAR(20) CHECK (avatar_source IN ('starter', 'spirit_animal'));

CREATE INDEX idx_users_avatar_key ON users (avatar_key) WHERE avatar_key IS NOT NULL;
