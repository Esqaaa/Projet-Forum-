package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
)

// Fonction de recherche 
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    queryParam := r.URL.Query().Get("q")
    categoryParam := r.URL.Query().Get("category")
    currentUserID := GetLoggedUserID(r)

    var topics []models.Topic
    categories := []string{"Sport", "Musique", "Automobile", "Aviation", "Sciences", "Informatique"}

    if queryParam != "" {
        sqlQuery := `
            SELECT t.id, t.title, t.content, t.status, t.created_at, t.category, u.username, t.author_id, t.image_url,
                   (SELECT COUNT(*) FROM user_pins WHERE user_id = ? AND topic_id = t.id) AS is_pinned_by_user
            FROM topics t
            JOIN users u ON t.author_id = u.id
            WHERE (t.title LIKE ? OR t.content LIKE ?)`

        var rows *sql.Rows
        var err error
        searchTerm := "%" + queryParam + "%"

        if categoryParam != "" {
            sqlQuery += " AND t.category = ? ORDER BY t.created_at DESC"
            rows, err = database.DB.Query(sqlQuery, currentUserID, searchTerm, searchTerm, categoryParam)
        } else {
            sqlQuery += " ORDER BY t.created_at DESC"
            rows, err = database.DB.Query(sqlQuery, currentUserID, searchTerm, searchTerm)
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
            var pinnedCount int 

            err = rows.Scan(
                &t.ID, &t.Title, &t.Content, &t.Status, &rawDate, &t.Category, &t.Author, &t.AuthorID, &imageURL,
                &pinnedCount, 
            )
            if err != nil {
                fmt.Println("Erreur scan recherche :", err)
                continue
            }

            if imageURL.Valid {
                t.ImageURL = imageURL.String
            } else {
                t.ImageURL = ""
            }

            t.IsPinnedByUser = pinnedCount > 0
            t.CreatedAt = string(rawDate)
            topics = append(topics, t)
        }
    }

    data := map[string]interface{}{
        "Query":            queryParam,
        "SelectedCategory": categoryParam,
        "Categories":       categories, 
        "Topics":           topics,
        "CurrentUserID":    currentUserID,
    }
    RenderTemplate(w, r, "search.html", data)
}