package handlers

import (
    "database/sql"
	"forum/database"
	"forum/models"
	"net/http"
    "fmt"
    "io"
    "os"
    "time"
)

func GetLoggedUserID(r *http.Request) int {
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
    userID := GetLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if r.Method == http.MethodGet {
        RenderTemplate(w, "create_topic.html", nil)
        return
    }

    if r.Method == http.MethodPost {
        r.ParseMultipartForm(5 << 20) 

        title := r.FormValue("title")
        content := r.FormValue("content")
        category := r.FormValue("category")
        if category == "" {
            category = "Général"
        }

        var imagePath string

        file, handler, err := r.FormFile("image")
        if err == nil { 
            defer file.Close()
            fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), handler.Filename)
            imagePath = "/static/uploads/" + fileName
            dst, err := os.Create("." + imagePath)
            if err != nil {
                http.Error(w, "Erreur stockage image", 500)
                return
            }
            defer dst.Close()
            io.Copy(dst, file)
        }

        query := "INSERT INTO topics (title, content, author_id, status, image_url, category) VALUES (?, ?, ?, ?, ?, ?)"
        _, err = database.DB.Exec(query, title, content, userID, "ouvert", imagePath, category)
        
        if err != nil {
            http.Error(w, "Erreur création : "+err.Error(), 500)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return 
    }
}

func ViewTopicHandler(w http.ResponseWriter, r *http.Request) {
    topicID := r.URL.Query().Get("id")
    currentUserID := GetLoggedUserID(r)

    var t models.Topic
    var rawDate []byte
    var imageURL sql.NullString

    query := `
        SELECT t.id, t.title, t.content, t.status, t.is_pinned, t.created_at, t.category, u.username, t.author_id, t.image_url
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
        &t.Category, 
        &t.Author,
        &t.AuthorID,
        &imageURL,
    )

    if err != nil {
        fmt.Println("Erreur ViewTopic:", err) 
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    
    if imageURL.Valid {
        t.ImageURL = imageURL.String
    } else {
        t.ImageURL = ""
    }

    t.CreatedAt = string(rawDate)
    t.Date = string(rawDate)

    rows, err := database.DB.Query(`
		SELECT m.id, m.content, m.created_at, u.username, m.author_id,
		       (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id) AS likes_count,
		       (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id AND user_id = ?) AS has_liked,
		       (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id) AS dislikes_count,
		       (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id AND user_id = ?) AS has_disliked
		FROM messages m 
		JOIN users u ON m.author_id = u.id 
		WHERE m.topic_id = ? 
		ORDER BY m.created_at ASC`, currentUserID, currentUserID, topicID)
    
    if err == nil {
        defer rows.Close()
    }
    
    var comments []models.Comment
    for rows != nil && rows.Next() {
        var c models.Comment
        var cDate []byte
        var hasLikedCount int
        var hasDislikedCount int  

        err = rows.Scan(&c.ID, &c.Content, &cDate, &c.Author, &c.AuthorID, &c.LikesCount, &hasLikedCount, &c.DislikesCount, &hasDislikedCount)
        if err != nil {
            fmt.Println("Erreur Scan message:", err)
            continue
        }

        c.HasLiked = hasLikedCount > 0
        c.HasDisliked = hasDislikedCount > 0
        c.Date = string(cDate)
        comments = append(comments, c)
    }

    data := map[string]interface{}{
        "Topic":         t,
        "Comments":      comments,
        "CurrentUserID": currentUserID,
    }
    RenderTemplate(w, "view_topic.html", data)
}

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    userID := GetLoggedUserID(r)
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

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
    userID := GetLoggedUserID(r)
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

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
    userID := GetLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    messageID := r.URL.Query().Get("id")
    topicID := r.URL.Query().Get("topic_id")

    var topicAuthorID int
    topicQuery := `SELECT t.author_id FROM topics t 
                   JOIN messages m ON t.id = m.topic_id 
                   WHERE m.id = ?`
    err := database.DB.QueryRow(topicQuery, messageID).Scan(&topicAuthorID)
    if err != nil {
        fmt.Println("Erreur recherche auteur topic:", err)
        http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
        return
    }

    var messageAuthorID int
    messageQuery := `SELECT author_id FROM messages WHERE id = ?`
    err = database.DB.QueryRow(messageQuery, messageID).Scan(&messageAuthorID)
    if err != nil {
        fmt.Println("Erreur recherche auteur message:", err)
        http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
        return
    }

    if userID == topicAuthorID || userID == messageAuthorID {
        _, err = database.DB.Exec("DELETE FROM messages WHERE id = ?", messageID)
        if err != nil {
            fmt.Println("Erreur lors de la suppression en BDD:", err)
        }
        http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
        return
    }

    http.Error(w, "Action non autorisée", http.StatusForbidden)
}

func EditMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    userID := GetLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    messageID := r.FormValue("message_id")
    topicID := r.FormValue("topic_id") 
    newContent := r.FormValue("content")

    if newContent == "" {
        http.Error(w, "Le contenu ne peut pas être vide", http.StatusBadRequest)
        return
    }

    var commentAuthorID int
    query := "SELECT author_id FROM messages WHERE id = ?"
    err := database.DB.QueryRow(query, messageID).Scan(&commentAuthorID)
    if err != nil {
        http.Error(w, "Message introuvable", http.StatusNotFound)
        return
    }

    if userID != commentAuthorID {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    _, err = database.DB.Exec("UPDATE messages SET content = ? WHERE id = ?", newContent, messageID)
    if err != nil {
        http.Error(w, "Erreur BDD", 500)
        return
    }

    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    userID := GetLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    messageID := r.FormValue("message_id")
    topicID := r.FormValue("topic_id")
    newContent := r.FormValue("content")

    if newContent == "" {
        http.Error(w, "Le contenu ne peut pas être vide", http.StatusBadRequest)
        return
    }

    var commentAuthorID int
    query := "SELECT author_id FROM messages WHERE id = ?"
    err := database.DB.QueryRow(query, messageID).Scan(&commentAuthorID)

    if err != nil {
        http.Error(w, "Message introuvable", http.StatusNotFound)
        return
    }

    if userID != commentAuthorID {
        http.Error(w, "Action non autorisée : vous n'êtes pas l'auteur de ce commentaire", http.StatusForbidden)
        return
    }

    _, err = database.DB.Exec("UPDATE messages SET content = ? WHERE id = ?", newContent, messageID)
    if err != nil {
        http.Error(w, "Erreur lors de la mise à jour du message : "+err.Error(), 500)
        return
    }

    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}

func UpdateTopicStatusHandler(w http.ResponseWriter, r *http.Request) {
    userID := GetLoggedUserID(r)
    topicID := r.URL.Query().Get("id")
    newStatus := r.URL.Query().Get("status") 

    var authorID int
    err := database.DB.QueryRow("SELECT author_id FROM topics WHERE id = ?", topicID).Scan(&authorID)
    
    if err != nil || userID != authorID {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    validStatuses := map[string]bool{"ouvert": true, "fermé": true, "archivé": true}
    if !validStatuses[newStatus] {
        http.Error(w, "Statut invalide", http.StatusBadRequest)
        return
    }

    _, err = database.DB.Exec("UPDATE topics SET status = ? WHERE id = ?", newStatus, topicID)
    if err != nil {
        http.Error(w, "Erreur BDD", 500)
        return
    }

    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}
 
func PinTopicHandler(w http.ResponseWriter, r *http.Request) {
    userID := GetLoggedUserID(r)
    topicID := r.URL.Query().Get("id")

    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    _, err := database.DB.Exec("UPDATE topics SET is_pinned = NOT is_pinned WHERE id = ?", topicID)
    if err != nil {
        fmt.Println("Erreur lors de l'UPDATE is_pinned:", err)
        http.Error(w, "Erreur lors de l'épinglage", 500)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}