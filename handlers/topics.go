package handlers

import (
	"forum/database"
	"forum/models"
	"net/http"
)

// getLoggedUserID : Récupère l'ID via le cookie "session"
func getLoggedUserID(r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0
	}

	var userID int
	// On cherche l'ID qui correspond au pseudo ou email stocké dans le cookie
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
		
		// INSERT sans image_url
		_, err := database.DB.Exec(
			"INSERT INTO topics (title, content, author_id, status) VALUES (?, ?, ?, ?)",
			title, content, userID, "ouvert",
		)
		if err != nil {
			http.Error(w, "Erreur création : "+err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ViewTopicHandler(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("id")

	var t models.Topic
	var rawDate []byte

	// 1. SELECT corrigé : On a enlevé t.image_url ici
	query := `
		SELECT t.id, t.title, t.content, t.status, t.created_at, u.username 
		FROM topics t 
		JOIN users u ON t.author_id = u.id 
		WHERE t.id = ?`
	
	// 2. SCAN corrigé : L'ordre doit être identique au SELECT ci-dessus
	err := database.DB.QueryRow(query, topicID).Scan(
		&t.ID, 
		&t.Title, 
		&t.Content, 
		&t.Status, 
		&rawDate, 
		&t.Author,
	)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// On remplit les deux champs de ton modèle pour être sûr
	t.CreatedAt = string(rawDate)
	t.Date = string(rawDate)

	// 3. Récupération des commentaires (messages)
	rows, err := database.DB.Query(`
		SELECT m.content, m.created_at, u.username 
		FROM messages m 
		JOIN users u ON m.author_id = u.id 
		WHERE m.topic_id = ? 
		ORDER BY m.created_at ASC`, topicID)
	
	if err != nil {
		// On continue même s'il n'y a pas de commentaires
	} else {
		defer rows.Close()
	}
	
	var comments []models.Comment
	for rows != nil && rows.Next() {
		var c models.Comment
		var cDate []byte
		rows.Scan(&c.Content, &cDate, &c.Author)
		c.Date = string(cDate)
		comments = append(comments, c)
	}

	data := map[string]interface{}{
		"Topic":    t,
		"Comments": comments,
	}
	RenderTemplate(w, "view_topic.html", data)
}