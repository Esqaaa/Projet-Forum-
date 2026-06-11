CREATE DATABASE IF NOT EXISTS forum_project;
USE forum_project;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    biography TEXT DEFAULT '',
    avatar_url VARCHAR(255) DEFAULT '/static/uploads/default-avatar.png',
    last_login DATETIME DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role ENUM('user', 'admin') DEFAULT 'user',
    is_banned BOOLEAN DEFAULT 0
);
