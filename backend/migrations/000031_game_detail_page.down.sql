ALTER TABLE games
    DROP COLUMN IF EXISTS catalog_hero_url,
    DROP COLUMN IF EXISTS how_to_play,
    DROP COLUMN IF EXISTS tutorial_url,
    DROP COLUMN IF EXISTS screenshots;
