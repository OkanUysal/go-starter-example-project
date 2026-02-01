-- Create example_token_blacklist table
CREATE TABLE IF NOT EXISTS example_token_blacklist (
    jti VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for cleanup queries
CREATE INDEX IF NOT EXISTS idx_example_token_blacklist_expires_at ON example_token_blacklist(expires_at);
CREATE INDEX IF NOT EXISTS idx_example_token_blacklist_user_id ON example_token_blacklist(user_id);
