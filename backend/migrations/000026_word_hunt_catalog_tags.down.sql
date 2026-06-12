UPDATE games
SET
    short_description = 'Teams race to guess words from one-word clues. Party word game.',
    tags = ARRAY['party', 'co-op', 'words']
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';
