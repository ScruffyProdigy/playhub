ALTER TABLE game_queues
    DROP COLUMN IF EXISTS forming_match_id,
    DROP COLUMN IF EXISTS party_id;

DROP TABLE IF EXISTS forming_match_assignments;
DROP TABLE IF EXISTS forming_matches;
DROP TABLE IF EXISTS party_members;
DROP TABLE IF EXISTS parties;
