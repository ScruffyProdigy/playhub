-- Catalog hero PNGs + rename Spock → Robot (slug + display name).
UPDATE games
SET
    slug = 'rock-paper-scissors-lizard-robot',
    name = 'Rock Paper Scissors Lizard Robot',
    hero_url = '/games/rpslr-hero.png'
WHERE slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET hero_url = '/games/word-hunt-hero.png'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
