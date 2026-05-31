DROP INDEX IF EXISTS idx_game_queues_one_waiting_per_user_mode_queue;

ALTER TABLE game_sessions
    DROP COLUMN IF EXISTS mode_queue_id,
    DROP COLUMN IF EXISTS mode_id;
