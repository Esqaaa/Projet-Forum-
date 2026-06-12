ALTER TABLE users ADD COLUMN biography TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255) DEFAULT '/static/uploads/default-avatar.png';
ALTER TABLE users ADD COLUMN last_login DATETIME DEFAULT CURRENT_TIMESTAMP;

