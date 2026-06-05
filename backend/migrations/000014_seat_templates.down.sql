UPDATE game_mode_seats
SET seat_key = 'a', sort_order = 0
WHERE mode_id = 'a2000000-0000-4000-8000-000000000001' AND seat_key = '1';

UPDATE game_mode_seats
SET seat_key = 'b', sort_order = 1
WHERE mode_id = 'a2000000-0000-4000-8000-000000000001' AND seat_key = '2';

UPDATE game_modes
SET seat_template = NULL
WHERE id = 'a2000000-0000-4000-8000-000000000001';

ALTER TABLE game_mode_seats
    DROP COLUMN IF EXISTS affinity_key,
    DROP COLUMN IF EXISTS queue_path;

ALTER TABLE game_modes
    DROP COLUMN IF EXISTS seat_template;
