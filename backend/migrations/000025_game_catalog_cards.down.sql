ALTER TABLE games
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS short_description,
    DROP COLUMN IF EXISTS icon_url;
