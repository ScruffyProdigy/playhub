-- Catalog hero banner for game list cards (required; separate from small icon_url).
ALTER TABLE games
    ADD COLUMN IF NOT EXISTS hero_url TEXT;

UPDATE games
SET hero_url = '/games/rpsls-hero.svg'
WHERE slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET hero_url = '/games/word-hunt-hero.svg'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';

UPDATE games
SET hero_url = '/games/default-hero.svg'
WHERE hero_url IS NULL OR trim(hero_url) = '';

ALTER TABLE games
    ALTER COLUMN hero_url SET NOT NULL,
    ALTER COLUMN hero_url SET DEFAULT '/games/default-hero.svg';
