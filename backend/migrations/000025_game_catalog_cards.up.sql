-- Catalog card metadata: icon (required), short blurb, tags for game list UI.
ALTER TABLE games
    ADD COLUMN IF NOT EXISTS icon_url TEXT,
    ADD COLUMN IF NOT EXISTS short_description TEXT,
    ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}';

UPDATE games
SET
    icon_url = '/games/rpsls.svg',
    short_description = 'Best-of-five duel with move cooldown. Fast reads, big throws.',
    tags = ARRAY['competitive', '1v1', 'quick']
WHERE slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET
    icon_url = '/games/word-hunt-board.svg',
    short_description = 'Everyone competes on a shared word grid. Push-your-luck party scoring.',
    tags = ARRAY['party', 'competitive', 'words']
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';

UPDATE games
SET icon_url = '/games/default.svg'
WHERE icon_url IS NULL OR trim(icon_url) = '';

ALTER TABLE games
    ALTER COLUMN icon_url SET NOT NULL,
    ALTER COLUMN icon_url SET DEFAULT '/games/default.svg';
