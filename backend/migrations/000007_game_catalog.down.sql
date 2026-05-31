ALTER TABLE game_queues DROP COLUMN IF EXISTS mode_queue_id;

DROP TABLE IF EXISTS mode_queues;
DROP TABLE IF EXISTS game_mode_seats;
DROP TABLE IF EXISTS game_modes;

DROP INDEX IF EXISTS idx_games_slug_unique;

ALTER TABLE games
    DROP COLUMN IF EXISTS webhook_secret,
    DROP COLUMN IF EXISTS game_version,
    DROP COLUMN IF EXISTS manifest_json,
    DROP COLUMN IF EXISTS manifest_synced_at,
    DROP COLUMN IF EXISTS manifest_etag,
    DROP COLUMN IF EXISTS manifest_hash;
