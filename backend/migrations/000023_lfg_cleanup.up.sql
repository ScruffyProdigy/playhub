-- LFG cleanup: drop deprecated party.together and unused forming match status.

ALTER TABLE parties DROP COLUMN IF EXISTS together;

ALTER TABLE forming_matches DROP CONSTRAINT IF EXISTS forming_matches_status_check;
ALTER TABLE forming_matches ADD CONSTRAINT forming_matches_status_check
    CHECK (status IN ('filling', 'fired'));
