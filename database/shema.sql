CREATE DATABASE IF NOT EXISTS forum_project;
USE forum_project;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    role ENUM('user', 'admin') DEFAULT 'user'
    is_banned BOOLEAN DEFAULT 0
);

-- Fichier SQL pour la mission FT-3 - Création d'un topic

-- Création de la table "topics"
CREATE TABLE IF NOT EXISTS topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    tags VARCHAR(255), 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
    author_id INT NOT NULL,
    is_pinned BOOLEAN DEFAULT 0,
    category TEXT NOT NULL DEFAULT 'Général',
    image_url TEXT, 
    status ENUM('ouvert', 'fermé', 'archivé') DEFAULT 'ouvert', 
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE 
);


-- Fichier SQL pour la mission FT-4 - Consulter le topic

-- Création de la table "messages"
CREATE TABLE IF NOT EXISTS messages (
    id int AUTO_INCREMENT PRIMARY KEY,
    topic_id INT NOT NULL,
    author_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Likes des messages
CREATE TABLE IF NOT EXISTS message_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_user_message_like (message_id, user_id),
    
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Dislikes des messages
CREATE TABLE IF NOT EXISTS message_dislikes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_user_message_dislike (message_id, user_id),
    
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE users ADD COLUMN biography TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255) DEFAULT '/static/uploads/default-avatar.png';
ALTER TABLE users ADD COLUMN last_login DATETIME DEFAULT CURRENT_TIMESTAMP;