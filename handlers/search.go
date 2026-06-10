package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	categoryParam := r.URL.Query().Get("category")
	currentUserID := GetLoggedUserID(r)

	var topics []models.Topic

	// On reprend la liste des catégories faites dans home.go pour les afficher dans le menu de recherche
	// Si il veut filtrer la recherche par catégorie
	categories := []string{"Sport", "Musique", "Automobile", "Aviation", "Sciences", "Informatique"}

	if queryParam != "" {
		// Base de la requête SQL
		sqlQuery := `
			SELECT t.id, t.title, t.content, t.status, t.is_pinned, t.created_at, t.category, u.username, t.author_id, t.image_url
			FROM topics t
			JOIN users u ON t.author_id = u.id
			WHERE (t.title LIKE ? OR t.content LIKE ?)`

		var rows *sql.Rows
		var err error
		searchTerm := "%" + queryParam + "%"

		// Si l'utilisateur a sélectionné une catégorie spécifique pour filtrer sa recherche
		if categoryParam != "" {
			sqlQuery += " AND t.category = ? ORDER BY t.is_pinned DESC, t.created_at DESC"
			rows, err = database.DB.Query(sqlQuery, searchTerm, searchTerm, categoryParam)
		} else {
			sqlQuery += " ORDER BY t.is_pinned DESC, t.created_at DESC"
			rows, err = database.DB.Query(sqlQuery, searchTerm, searchTerm)
		}

		if err != nil {
			fmt.Println("Erreur lors de la recherche :", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var t models.Topic
			var rawDate []byte
			var imageURL sql.NullString

			err = rows.Scan(&t.ID, &t.Title, &t.Content, &t.Status, &t.IsPinned, &rawDate, &t.Category, &t.Author, &t.AuthorID, &imageURL)
			if err != nil {
				fmt.Println("Erreur scan recherche :", err)
				continue
			}

			if imageURL.Valid {
				t.ImageURL = imageURL.String
			}
			t.CreatedAt = string(rawDate)
			topics = append(topics, t)
		}
	}

	// On envoie les données au HTML
	data := map[string]interface{}{
		"Query":            queryParam,
		"SelectedCategory": categoryParam,
		"Categories":       categories, 
		"Topics":           topics,
		"CurrentUserID":    currentUserID,
	}
	RenderTemplate(w, "search.html", data)
}