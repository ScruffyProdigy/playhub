-- One waiting queue row per user across all games (not per game).
DROP INDEX IF EXISTS idx_game_queues_one_waiting_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS idx_game_queues_one_waiting_per_user
    ON game_queues (user_id)
    WHERE status = 'waiting';
