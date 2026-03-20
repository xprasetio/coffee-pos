-- +migrate Down
-- Rename email back to username
ALTER TABLE users CHANGE COLUMN email username VARCHAR(50) UNIQUE NOT NULL;

-- Rename name back to full_name
ALTER TABLE users CHANGE COLUMN name full_name VARCHAR(100) NOT NULL;

-- Rename password back to password_hash
ALTER TABLE users CHANGE COLUMN password password_hash VARCHAR(255) NOT NULL;

-- Update index back to idx_username
DROP INDEX idx_email ON users;
CREATE INDEX idx_username ON users(username);
