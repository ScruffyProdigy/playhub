-- Social chat rooms (Step 1: no tables yet).

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invite_code VARCHAR(8) NOT NULL,
    host_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT rooms_status_check CHECK (status IN ('open', 'closed'))
);

CREATE UNIQUE INDEX idx_rooms_invite_code ON rooms (UPPER(invite_code));

CREATE TABLE room_members (
    room_id UUID NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

-- One active room membership per player globally.
CREATE UNIQUE INDEX idx_room_members_one_room_per_user ON room_members (user_id);

CREATE TABLE room_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT room_messages_body_length CHECK (char_length(body) BETWEEN 1 AND 2000)
);

CREATE INDEX idx_room_messages_room_created ON room_messages (room_id, created_at DESC);
