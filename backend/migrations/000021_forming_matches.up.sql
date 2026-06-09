-- LFG Phase B: parties, forming matches, and queue linkage.

CREATE TABLE parties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    leader_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    mode_queue_id UUID NOT NULL REFERENCES mode_queues(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    together BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT parties_status_check CHECK (status IN ('waiting', 'placed', 'matched', 'cancelled'))
);

CREATE INDEX idx_parties_mode_queue_status ON parties (mode_queue_id, status);

CREATE TABLE party_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    queue_path VARCHAR(100) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (party_id, user_id)
);

CREATE INDEX idx_party_members_user ON party_members (user_id);

CREATE TABLE forming_matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mode_queue_id UUID NOT NULL REFERENCES mode_queues(id) ON DELETE CASCADE,
    mode_id UUID NOT NULL REFERENCES game_modes(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'filling',
    target_spec JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fired_at TIMESTAMPTZ,
    CONSTRAINT forming_matches_status_check CHECK (status IN ('filling', 'ready', 'fired'))
);

CREATE UNIQUE INDEX idx_forming_matches_one_filling_per_queue
    ON forming_matches (mode_queue_id)
    WHERE status = 'filling';

CREATE INDEX idx_forming_matches_mode_queue ON forming_matches (mode_queue_id, status);

CREATE TABLE forming_match_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    forming_match_id UUID NOT NULL REFERENCES forming_matches(id) ON DELETE CASCADE,
    seat_key VARCHAR(200) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    party_id UUID REFERENCES parties(id) ON DELETE SET NULL,
    queue_path VARCHAR(100) NOT NULL DEFAULT '',
    affinity_key VARCHAR(100),
    source VARCHAR(20) NOT NULL DEFAULT 'solo',
    table_id UUID REFERENCES room_tables(id) ON DELETE SET NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (forming_match_id, seat_key),
    CONSTRAINT forming_match_assignments_source_check CHECK (source IN ('party', 'solo', 'table'))
);

CREATE INDEX idx_forming_match_assignments_user
    ON forming_match_assignments (forming_match_id, user_id)
    WHERE user_id IS NOT NULL;

ALTER TABLE game_queues
    ADD COLUMN IF NOT EXISTS party_id UUID REFERENCES parties(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS forming_match_id UUID REFERENCES forming_matches(id) ON DELETE SET NULL;

CREATE INDEX idx_game_queues_forming_match ON game_queues (forming_match_id)
    WHERE forming_match_id IS NOT NULL;

CREATE INDEX idx_game_queues_party ON game_queues (party_id)
    WHERE party_id IS NOT NULL;
