package handlers

import (
	"database/sql"
	"net/http"
	"time"
)

// FT-3 : Création d'un topic
func CreateTopicHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.ServeFile(w, r, "./templates/html/create_topic.html")
			return
		}

		// Récupérer les données du formulaire
		title := r.FormValue("title")
		content := r.FormValue("content")
		tags := r.FormValue("tags")

		// Remplacer le 1 par l'ID de l'utilisateur connecté (ex: depuis la session)
		// Normalement dans la fonction de session, on aurait une fonction GetUserIDFromSession(r) qui retourne l'ID de l'utilisateur connecté
		authorID := 1 

		// Insertion en base de données
		_, err := db.Exec(`INSERT INTO topics (title, content, tags, author_id, status, created_at) 
						   VALUES (?, ?, ?, ?, 'ouvert', ?)`,
			title, content, tags, authorID, time.Now())

		if err != nil {
			http.Error(w, "Erreur lors de la création du topic", http.StatusInternalServerError)
			return
		}

		// Redirection vers l'accueil ou le topic
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}