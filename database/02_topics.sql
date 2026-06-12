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