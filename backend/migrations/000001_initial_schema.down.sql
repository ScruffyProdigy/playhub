-- Rollback for initial schema migration
-- This will drop all tables and functions created in the up migration

-- Drop triggers first
DROP TRIGGER IF EXISTS update_digital_goods_updated_at ON digital_goods;
DROP TRIGGER IF EXISTS update_games_updated_at ON games;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS user_inventory;
DROP TABLE IF EXISTS digital_goods;
DROP TABLE IF EXISTS game_session_participants;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS game_queues;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS magic_links;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
