-- +migrate Up
-- Rename username to email
ALTER TABLE users CHANGE COLUMN username email VARCHAR(50) UNIQUE NOT NULL;

-- Rename full_name to name
ALTER TABLE users CHANGE COLUMN full_name name VARCHAR(100) NOT NULL;

-- Rename password_hash to password
ALTER TABLE users CHANGE COLUMN password_hash password VARCHAR(255) NOT NULL;

-- Update index from idx_username to idx_email
DROP INDEX idx_username ON users;
CREATE INDEX idx_email ON users(email);
