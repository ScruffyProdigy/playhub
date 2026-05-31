-- All playable games must have catalog modes and mode queues. Seed the RPS demo; retire Party Lobby.

UPDATE games
SET
    category = 'catalog',
    status = 'active',
    manifest_hash = COALESCE(manifest_hash, 'seed-rps-duel-v1'),
    manifest_synced_at = COALESCE(manifest_synced_at, NOW()),
    game_version = COALESCE(game_version, '1.0.0')
WHERE id = 'a1000000-0000-4000-8000-000000000001';

UPDATE games
SET status = 'inactive'
WHERE id = 'a1000000-0000-4000-8000-000000000002';

INSERT INTO game_modes (id, game_id, mode_key, display_name, min_players, max_players, best_of, status)
VALUES (
    'a2000000-0000-4000-8000-000000000001',
    'a1000000-0000-4000-8000-000000000001',
    'duel',
    '1v1 Duel (best 3 of 5)',
    2,
    2,
    5,
    'active'
)
ON CONFLICT (game_id, mode_key) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    min_players = EXCLUDED.min_players,
    max_players = EXCLUDED.max_players,
    best_of = EXCLUDED.best_of,
    status = 'active',
    updated_at = NOW();

INSERT INTO game_mode_seats (mode_id, seat_key, sort_order)
VALUES
    ('a2000000-0000-4000-8000-000000000001', 'a', 0),
    ('a2000000-0000-4000-8000-000000000001', 'b', 1)
ON CONFLICT (mode_id, seat_key) DO NOTHING;

INSERT INTO mode_queues (id, mode_id, name, players_to_start, status, is_default)
VALUES (
    'a3000000-0000-4000-8000-000000000001',
    'a2000000-0000-4000-8000-000000000001',
    'Default',
    2,
    'active',
    true
)
ON CONFLICT (mode_id, name) DO UPDATE SET
    players_to_start = EXCLUDED.players_to_start,
    status = 'active',
    is_default = true,
    updated_at = NOW();
