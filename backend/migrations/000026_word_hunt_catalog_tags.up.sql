-- Word Hunt is individual competitive (not team vs team or co-op).
UPDATE games
SET
    short_description = 'Everyone competes on a shared word grid. Push-your-luck party scoring.',
    tags = ARRAY['party', 'competitive', 'words']
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
