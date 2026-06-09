-- Party tree structure (mirrors seatTemplate branches).

ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS party_tree JSONB NOT NULL DEFAULT '{}'::jsonb;
