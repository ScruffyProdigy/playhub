ALTER TABLE games ADD COLUMN IF NOT EXISTS play_url TEXT;

UPDATE games SET play_url = api_base_url WHERE play_url IS NULL AND api_base_url IS NOT NULL;
