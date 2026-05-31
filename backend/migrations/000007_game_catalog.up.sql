-- Game catalog: cached manifests, modes, seats, and matchmaking queues (mode_queues).
ALTER TABLE games
    ADD COLUMN IF NOT EXISTS manifest_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS manifest_etag TEXT,
    ADD COLUMN IF NOT EXISTS manifest_synced_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS manifest_json JSONB,
    ADD COLUMN IF NOT EXISTS game_version VARCHAR(100),
    ADD COLUMN IF NOT EXISTS webhook_secret VARCHAR(128);

CREATE UNIQUE INDEX IF NOT EXISTS idx_games_slug_unique ON games (slug) WHERE slug IS NOT NULL;

CREATE TABLE IF NOT EXISTS game_modes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mode_key VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    min_players INTEGER NOT NULL DEFAULT 1,
    max_players INTEGER NOT NULL DEFAULT 2,
    best_of INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (game_id, mode_key)
);

CREATE TABLE IF NOT EXISTS game_mode_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mode_id UUID NOT NULL REFERENCES game_modes(id) ON DELETE CASCADE,
    seat_key VARCHAR(100) NOT NULL,
    team VARCHAR(100),
    role VARCHAR(100),
    sort_order INTEGER NOT NULL DEFAULT 0,
    UNIQUE (mode_id, seat_key)
);

CREATE INDEX IF NOT EXISTS idx_game_mode_seats_mode_id ON game_mode_seats (mode_id, sort_order);

-- Matchmaking buckets (distinct from game_queues, which holds waiting players).
CREATE TABLE IF NOT EXISTS mode_queues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mode_id UUID NOT NULL REFERENCES game_modes(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    players_to_start INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (mode_id, name)
);

CREATE INDEX IF NOT EXISTS idx_mode_queues_mode_id ON mode_queues (mode_id);

ALTER TABLE game_queues
    ADD COLUMN IF NOT EXISTS mode_queue_id UUID REFERENCES mode_queues(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_game_queues_mode_queue_id ON game_queues (mode_queue_id)
    WHERE status = 'waiting';
