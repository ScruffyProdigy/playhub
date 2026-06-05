DROP INDEX IF EXISTS idx_game_queues_mode_queue_path_waiting;

ALTER TABLE game_queues
    DROP COLUMN IF EXISTS queue_path;
