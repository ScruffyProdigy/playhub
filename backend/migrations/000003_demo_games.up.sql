-- Demo games for local development and stakeholder demos.
INSERT INTO games (id, name, description, min_players, max_players, category, status)
VALUES
    (
        'a1000000-0000-4000-8000-000000000001',
        'Quick Match',
        'Jump in for a fast demo session with another player.',
        2,
        4,
        'demo',
        'active'
    ),
    (
        'a1000000-0000-4000-8000-000000000002',
        'Party Lobby',
        'Gather a group and start a session together.',
        2,
        8,
        'demo',
        'active'
    )
ON CONFLICT (id) DO NOTHING;
