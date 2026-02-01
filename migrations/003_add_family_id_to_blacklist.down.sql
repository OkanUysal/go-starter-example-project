-- Remove family_id column from example_token_blacklist table
DROP INDEX IF EXISTS idx_example_token_blacklist_family_id;
ALTER TABLE example_token_blacklist DROP COLUMN IF EXISTS family_id;
