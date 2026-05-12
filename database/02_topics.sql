-- Fichier SQL pour la mission FT-3 - Création d'un topic

-- Création de la table "topics"
CREATE TABLE IF NOT EXISTS topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    tags VARCHAR(255), 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
    author_id INT NOT NULL, 
    status ENUM('ouvert', 'fermé', 'archivé') DEFAULT 'ouvert', 
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE 
);


ALTER TABLE topics ADD COLUMN is_pinned BOOLEAN DEFAULT 0;
ALTER TABLE topics ADD COLUMN image_url TEXT;