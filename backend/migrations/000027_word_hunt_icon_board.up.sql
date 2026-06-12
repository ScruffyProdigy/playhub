-- Bust CDN cache for Word Hunt catalog icon (grid board replaces magnifying-glass art).
UPDATE games
SET icon_url = '/games/word-hunt-board.svg'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
