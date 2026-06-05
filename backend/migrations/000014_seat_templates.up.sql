-- Phase A: seatTemplate expansion stored on modes; expanded leaves carry affinity/queue metadata.

ALTER TABLE game_modes
    ADD COLUMN IF NOT EXISTS seat_template JSONB;

ALTER TABLE game_mode_seats
    ADD COLUMN IF NOT EXISTS affinity_key VARCHAR(100),
    ADD COLUMN IF NOT EXISTS queue_path VARCHAR(100);

UPDATE game_modes
SET seat_template = '{"count": 2}'::jsonb
WHERE id = 'a2000000-0000-4000-8000-000000000001';

UPDATE game_mode_seats
SET seat_key = '1', sort_order = 0
WHERE mode_id = 'a2000000-0000-4000-8000-000000000001' AND seat_key = 'a';

UPDATE game_mode_seats
SET seat_key = '2', sort_order = 1
WHERE mode_id = 'a2000000-0000-4000-8000-000000000001' AND seat_key = 'b';
