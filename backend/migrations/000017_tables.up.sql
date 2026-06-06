-- Forming game tables inside chat rooms (Step 2).

CREATE TABLE room_tables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    mode_id UUID NOT NULL REFERENCES game_modes (id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'forming',
    session_id UUID REFERENCES game_sessions (id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT room_tables_status_check CHECK (status IN ('forming', 'started', 'discarded'))
);

CREATE INDEX idx_room_tables_room_id ON room_tables (room_id);
CREATE INDEX idx_room_tables_room_forming ON room_tables (room_id) WHERE status = 'forming';

CREATE TABLE table_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_id UUID NOT NULL REFERENCES room_tables (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    seat_key VARCHAR(128) NOT NULL,
    seated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (table_id, seat_key),
    UNIQUE (user_id)
);

CREATE INDEX idx_table_seats_table_id ON table_seats (table_id);
