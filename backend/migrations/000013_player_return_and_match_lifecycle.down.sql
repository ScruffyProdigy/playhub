DROP INDEX IF EXISTS idx_game_session_participants_session_user;

ALTER TABLE game_session_participants
    DROP COLUMN IF EXISTS finished_at,
    DROP COLUMN IF EXISTS return_context;
