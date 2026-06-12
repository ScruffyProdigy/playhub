-- Catalog icons: dedicated PNG artwork for Word Hunt and RPSLR.
UPDATE games
SET icon_url = '/games/rpslr-icon.png'
WHERE slug = 'rock-paper-scissors-lizard-robot'
   OR slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET icon_url = '/games/word-hunt-icon.png'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
