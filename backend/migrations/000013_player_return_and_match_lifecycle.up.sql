-- Per-player return routing after a match and lifecycle timestamps for game callbacks.

ALTER TABLE game_session_participants
    ADD COLUMN IF NOT EXISTS return_context JSONB,
    ADD COLUMN IF NOT EXISTS finished_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_game_session_participants_session_user
    ON game_session_participants (session_id, user_id);
