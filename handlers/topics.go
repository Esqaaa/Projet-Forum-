package handlers

import (
	"forum/database"
	"forum/models"
	"net/http"
    "fmt"
)

func getLoggedUserID(r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0
	}

	var userID int
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

    query := `
        SELECT t.id, t.title, t.content, t.status, t.is_pinned, t.created_at, u.username, t.author_id
        FROM topics t
        JOIN users u ON t.author_id = u.id
        WHERE t.id = ?`

    err := database.DB.QueryRow(query, topicID).Scan(
        &t.ID, 
        &t.Title, 
        &t.Content, 
        &t.Status, 
        &t.IsPinned, 
        &rawDate, 
        &t.Author,
        &t.AuthorID,
    )

    if err != nil {
        fmt.Println("Erreur ViewTopic:", err) 
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    
    t.CreatedAt = string(rawDate)
    t.Date = string(rawDate)

    rows, err := database.DB.Query(`
        SELECT m.id, m.content, m.created_at, u.username 
        FROM messages m 
        JOIN users u ON m.author_id = u.id 
        WHERE m.topic_id = ? 
        ORDER BY m.created_at ASC`, topicID)
    
    if err == nil {
        defer rows.Close()
    }
    
    var comments []models.Comment
    for rows != nil && rows.Next() {
        var c models.Comment
        var cDate []byte
        rows.Scan(&c.ID, &c.Content, &cDate, &c.Author)
        c.Date = string(cDate)
        comments = append(comments, c)
    }

    data := map[string]interface{}{
        "Topic":         t,
        "Comments":      comments,
        "CurrentUserID": getLoggedUserID(r),
    }
    RenderTemplate(w, "view_topic.html", data)
}

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    userID := getLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    topicID := r.FormValue("topic_id")
    content := r.FormValue("content")

    if content != "" {
        _, err := database.DB.Exec(
            "INSERT INTO messages (topic_id, author_id, content) VALUES (?, ?, ?)",
            topicID, userID, content,
        )
        if err != nil {
            http.Error(w, "Erreur lors de l'envoi du message : "+err.Error(), 500)
            return
        }
    }

    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}


// DeleteTopicHandler supprime un topic et ses messages (grâce au ON DELETE CASCADE en SQL)
func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
    userID := getLoggedUserID(r)
    topicID := r.URL.Query().Get("id")
    var authorID int
    
    err := database.DB.QueryRow("SELECT author_id FROM topics WHERE id = ?", topicID).Scan(&authorID)

    if err != nil || userID != authorID {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    database.DB.Exec("DELETE FROM topics WHERE id = ?", topicID)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteMessageHandler permet au proprio du topic de supprimer un message
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
    userID := getLoggedUserID(r)
    messageID := r.URL.Query().Get("id")
    topicID := r.URL.Query().Get("topic_id")

    // Vérifier si l'user est le proprio du TOPIC
    var topicAuthorID int
    query := `SELECT t.author_id FROM topics t 
              JOIN messages m ON t.id = m.topic_id 
              WHERE m.id = ?`
    err := database.DB.QueryRow(query, messageID).Scan(&topicAuthorID)

    if err != nil || userID != topicAuthorID {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    database.DB.Exec("DELETE FROM messages WHERE id = ?", messageID)
    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}

// Ajouter un topic aux favoris 
func PinTopicHandler(w http.ResponseWriter, r *http.Request) {
	userID := getLoggedUserID(r)
	topicID := r.URL.Query().Get("id")

	// Si l'utilisateur n'est pas connecté, on redirige
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err := database.DB.Exec("UPDATE topics SET is_pinned = NOT is_pinned WHERE id = ?", topicID)
	if err != nil {
		fmt.Println("Erreur SQL Pin:", err)
		http.Error(w, "Erreur lors de l'épinglage", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}