ALTER TABLE avatar_readings
    ADD COLUMN phase_started_at TIMESTAMPTZ;

UPDATE avatar_readings
SET phase_started_at = updated_at
WHERE status IN ('generating_questions', 'processing');
