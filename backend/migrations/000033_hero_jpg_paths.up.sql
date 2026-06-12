-- Point catalog/detail heroes at JPEG assets (cache-bust from legacy .png / .svg paths).
UPDATE games
SET hero_url = '/games/word-hunt-hero.jpg'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';

UPDATE games
SET hero_url = '/games/rpslr-hero.jpg'
WHERE slug = 'rock-paper-scissors-lizard-robot'
   OR slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000002';
