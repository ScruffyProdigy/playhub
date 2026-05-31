-- Mode-queue scoped matchmaking: session metadata and one waiting row per user per queue.
ALTER TABLE game_sessions
    ADD COLUMN IF NOT EXISTS mode_id UUID REFERENCES game_modes(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS mode_queue_id UUID REFERENCES mode_queues(id) ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_game_queues_one_waiting_per_user_mode_queue
    ON game_queues (mode_queue_id, user_id)
    WHERE status = 'waiting' AND mode_queue_id IS NOT NULL;
