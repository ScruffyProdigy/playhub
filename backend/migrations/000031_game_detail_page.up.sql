-- Dedicated game pages: optional catalog hero, long copy, tutorial, screenshots.
ALTER TABLE games
    ADD COLUMN IF NOT EXISTS catalog_hero_url TEXT,
    ADD COLUMN IF NOT EXISTS how_to_play TEXT,
    ADD COLUMN IF NOT EXISTS tutorial_url TEXT,
    ADD COLUMN IF NOT EXISTS screenshots TEXT[] NOT NULL DEFAULT '{}';

-- Word Hunt detail page content (description column = long form).
UPDATE games
SET
    description = 'Word Hunt is a competitive party word game on a shared grid. Clue-givers give one-word hints; everyone else races to guess words tile by tile. The pot grows while you wait — pass early for the biggest score, or push your luck for one more guess.',
    how_to_play = 'Each round a clue-giver picks a one-word clue tied to words on the board. Other players guess which tiles match; correct guesses score from a shared pot that grows over time. Wrong guesses can crash the pot for that clue. First color to clear all their words wins the match.',
    tutorial_url = 'https://word-hunt-arena.win'
WHERE slug = 'word-hunt' OR name ILIKE 'word hunt%';

UPDATE games
SET
    description = 'Rock Paper Scissors Lizard Robot is a best-of-five duel with a move cooldown twist. Standard RPSLR rules apply, but each move sits out for a few rounds after you play it — so you cannot spam the same throw.',
    how_to_play = 'Pick rock, paper, scissors, lizard, or robot each round. Each move beats two others (rock crushes scissors and lizard, paper covers rock and disproves robot, etc.). After you play a move it gains delay marks and cannot be chosen again until they tick down. Lizard and robot start on cooldown. First to three round wins takes the match.',
    tutorial_url = 'https://rpsls-duel.win'
WHERE slug = 'rock-paper-scissors-lizard-robot'
   OR slug = 'rock-paper-scissors-lizard-spock'
   OR id = 'a1000000-0000-4000-8000-000000000002';
