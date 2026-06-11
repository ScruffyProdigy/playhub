-- Game-minted launch URL bases (without JWT) persisted at match finalize for refresh paths.

ALTER TABLE game_session_participants
    ADD COLUMN IF NOT EXISTS launch_url_base TEXT;
