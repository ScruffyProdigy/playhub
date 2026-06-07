CREATE TABLE avatar_readings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL DEFAULT 'generating_questions',
    draw INT[] NOT NULL CHECK (array_length(draw, 1) = 5),
    questions_json JSONB,
    user_answers TEXT[],
    personality_json JSONB,
    totems_json JSONB,
    ranking_json JSONB,
    selected_totem_name VARCHAR(100),
    art_direction_version INT NOT NULL DEFAULT 1,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_avatar_readings_user_id ON avatar_readings (user_id);
CREATE INDEX idx_avatar_readings_status ON avatar_readings (status)
    WHERE completed_at IS NULL;

CREATE TABLE avatar_renders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reading_id UUID NOT NULL REFERENCES avatar_readings (id) ON DELETE CASCADE,
    totem_name VARCHAR(100) NOT NULL,
    art_direction_version INT NOT NULL,
    image_url TEXT NOT NULL,
    image_prompt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_avatar_renders_reading_id ON avatar_renders (reading_id);

ALTER TABLE users
    ADD COLUMN avatar_reading_id UUID REFERENCES avatar_readings (id);

CREATE INDEX idx_users_avatar_reading_id ON users (avatar_reading_id)
    WHERE avatar_reading_id IS NOT NULL;
