-- Create example_user table
CREATE TABLE IF NOT EXISTS example_user (
    id VARCHAR(255) PRIMARY KEY,
    guest_id VARCHAR(255) UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    is_guest BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_example_user_guest_id ON example_user(guest_id);
CREATE INDEX IF NOT EXISTS idx_example_user_google_id ON example_user(google_id);
CREATE INDEX IF NOT EXISTS idx_example_user_role ON example_user(role);

-- Add check constraint for role
ALTER TABLE example_user 
ADD CONSTRAINT chk_user_role CHECK (role IN ('USER', 'ADMIN'));
