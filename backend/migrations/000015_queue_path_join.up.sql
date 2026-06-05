-- Persist player's composition queue path on waiting rows.

ALTER TABLE game_queues
    ADD COLUMN IF NOT EXISTS queue_path VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_game_queues_mode_queue_path_waiting
    ON game_queues (mode_queue_id, queue_path)
    WHERE status = 'waiting';
