-- Add family_id column to example_token_blacklist table
ALTER TABLE example_token_blacklist ADD COLUMN IF NOT EXISTS family_id VARCHAR(255);

-- Create index for family_id lookups
CREATE INDEX IF NOT EXISTS idx_example_token_blacklist_family_id ON example_token_blacklist(family_id);
