UPDATE games
SET icon_url = '/games/rpsls.svg'
WHERE slug IN ('rock-paper-scissors-lizard-robot', 'rock-paper-scissors-lizard-spock')
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET icon_url = '/games/word-hunt-board.svg'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
