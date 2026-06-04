DROP INDEX IF EXISTS idx_game_queues_one_waiting_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS idx_game_queues_one_waiting_per_user
    ON game_queues (game_id, user_id)
    WHERE status = 'waiting';
