ALTER TABLE forming_matches DROP CONSTRAINT IF EXISTS forming_matches_status_check;
ALTER TABLE forming_matches ADD CONSTRAINT forming_matches_status_check
    CHECK (status IN ('filling', 'ready', 'fired'));

ALTER TABLE parties ADD COLUMN IF NOT EXISTS together BOOLEAN NOT NULL DEFAULT false;
