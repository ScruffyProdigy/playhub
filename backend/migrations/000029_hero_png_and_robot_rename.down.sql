UPDATE games
SET
    slug = 'rock-paper-scissors-lizard-spock',
    name = 'Rock Paper Scissors Lizard Spock',
    hero_url = '/games/rpsls-hero.svg'
WHERE slug = 'rock-paper-scissors-lizard-robot'
   OR id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET hero_url = '/games/word-hunt-hero.svg'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
