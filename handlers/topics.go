package handlers

import (
	"forum/database"
	"net/http"
)

// getLoggedUserID : Traduit le cookie de ton pote en ID utilisateur
func getLoggedUserID(r *http.Request) int {
	cookie, err := r.Cookie("session") // Nom exact utilisé dans auth.go
	if err != nil {
		return 0
	}

	var userID int
	// On cherche l'ID qui correspond au pseudo/email stocké dans le cookie
	query := "SELECT id FROM users WHERE username = ? OR email = ?"
	err = database.DB.QueryRow(query, cookie.Value, cookie.Value).Scan(&userID)
	if err != nil {
		return 0
	}
	return userID
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	userID := getLoggedUserID(r)
	if userID == 0 {
		// Si pas de cookie valide, on renvoie au login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		RenderTemplate(w, "create_topic.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		tags := r.FormValue("tags")

		_, err := database.DB.Exec(
			"INSERT INTO topics (title, content, tags, author_id) VALUES (?, ?, ?, ?)",
			title, content, tags, userID,
		)
		if err != nil {
			http.Error(w, "Erreur lors de la création : "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}