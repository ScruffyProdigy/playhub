-- Third-party game handoff: browser play URL + server-to-server API base.
ALTER TABLE games
    ADD COLUMN IF NOT EXISTS slug VARCHAR(100),
    ADD COLUMN IF NOT EXISTS play_url TEXT,
    ADD COLUMN IF NOT EXISTS api_base_url TEXT,
    ADD COLUMN IF NOT EXISTS game_mode VARCHAR(50) DEFAULT 'duel';

-- Demo quick match → local RPS game (see demo-game-rps/docs/lobby-protocol-handoff.md).
UPDATE games
SET
    slug = 'rock-paper-scissors-lizard-spock',
    name = 'Rock Paper Scissors Lizard Spock',
    description = 'Best-of duel with move cooldown. Demo third-party game.',
    play_url = 'http://localhost:5174',
    api_base_url = 'http://localhost:3001',
    game_mode = 'duel',
    min_players = 2,
    max_players = 2
WHERE id = 'a1000000-0000-4000-8000-000000000001';
