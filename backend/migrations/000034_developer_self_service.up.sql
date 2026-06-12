ALTER TABLE games ADD COLUMN IF NOT EXISTS owner_user_id UUID REFERENCES users(id);
ALTER TABLE games ADD COLUMN IF NOT EXISTS visibility TEXT;
UPDATE games SET visibility = 'public' WHERE visibility IS NULL;
ALTER TABLE games ALTER COLUMN visibility SET NOT NULL;
ALTER TABLE games ALTER COLUMN visibility SET DEFAULT 'draft';
ALTER TABLE games ADD CONSTRAINT games_visibility_check
  CHECK (visibility IN ('draft', 'private_testing', 'pending_review', 'public'));

ALTER TABLE games ADD COLUMN IF NOT EXISTS contact_email TEXT;
ALTER TABLE games ADD COLUMN IF NOT EXISTS website_url TEXT;
ALTER TABLE games ADD COLUMN IF NOT EXISTS community_url TEXT;

CREATE TABLE IF NOT EXISTS game_integration_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    check_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pass', 'fail', 'skipped')),
    message TEXT,
    detail_json JSONB,
    ran_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (game_id, check_id)
);

CREATE INDEX IF NOT EXISTS idx_game_integration_checks_game_id ON game_integration_checks(game_id);
