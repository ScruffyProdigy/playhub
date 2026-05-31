-- Legacy per-game matchmaking fields superseded by game_modes catalog.
ALTER TABLE games
    DROP COLUMN IF EXISTS game_mode,
    DROP COLUMN IF EXISTS min_players,
    DROP COLUMN IF EXISTS max_players;

ALTER TABLE game_modes
    DROP COLUMN IF EXISTS best_of;
