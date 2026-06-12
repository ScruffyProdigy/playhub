DROP TABLE IF EXISTS game_integration_checks;

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_visibility_check;
ALTER TABLE games DROP COLUMN IF EXISTS community_url;
ALTER TABLE games DROP COLUMN IF EXISTS website_url;
ALTER TABLE games DROP COLUMN IF EXISTS contact_email;
ALTER TABLE games DROP COLUMN IF EXISTS visibility;
ALTER TABLE games DROP COLUMN IF EXISTS owner_user_id;
