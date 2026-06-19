package handlers

import (
    "database/sql"
    "fmt"
    "forum/database"
    "forum/models"
    "io"
    "net/http"
    "os"
    "time"
)

func GetLoggedUserID(r *http.Request) int {
    cookie, err := r.Cookie("session")
    if err != nil {
        return 0
    }

	var userID int
	query := "SELECT id FROM users WHERE BINARY username = ? OR BINARY email = ?"
	err = database.DB.QueryRow(query, cookie.Value, cookie.Value).Scan(&userID)
	if err != nil {
		return 0
	}
	return userID
}

func GetLoggedUser(r *http.Request) (models.User, error) {
    cookie, err := r.Cookie("session")
    if err != nil {
        return models.User{}, err
    }

    var u models.User

    query := `
        SELECT id, username, email, role
        FROM users
        WHERE username = ? OR email = ?
    `

    err = database.DB.QueryRow(query, cookie.Value, cookie.Value).Scan(
        &u.ID,
        &u.Username,
        &u.Email,
        &u.Role,
    )

    if err != nil {
        return models.User{}, err
    }

    return u, nil
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
    userID := GetLoggedUserID(r)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if r.Method == http.MethodGet {
        RenderTemplate(w, r, "create_topic.html", nil)
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

    user, _ := GetLoggedUser(r)

    sortParam := r.URL.Query().Get("sort")
    if sortParam == "" {
        sortParam = "chrono" // Par défaut : plus récents au plus anciens
	}

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
        &t.IsPinnedByUser,
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

    isOwner := currentUserID != 0 && currentUserID == t.AuthorID

    if t.Status == "archivé" && currentUserID != t.AuthorID {
        http.Error(w, "Ce sujet a été archivé et n'est plus accessible.", http.StatusForbidden)
        return
    }
    
    if imageURL.Valid {
        t.ImageURL = imageURL.String
    } else {
        t.ImageURL = ""
    }

    t.CreatedAt = string(rawDate)
    t.Date = string(rawDate)

    var commentsQuery string
    if sortParam == "popular" {
        commentsQuery = `
            SELECT m.id, m.content, m.created_at, u.username, m.author_id,
                   (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id) AS likes_count,
                   (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id AND user_id = ?) AS has_liked,
                   (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id) AS dislikes_count,
                   (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id AND user_id = ?) AS has_disliked
            FROM messages m 
            JOIN users u ON m.author_id = u.id 
            WHERE m.topic_id = ? 
            ORDER BY (likes_count - dislikes_count) DESC, m.created_at DESC`
    } else {
        commentsQuery = `
            SELECT m.id, m.content, m.created_at, u.username, m.author_id,
                   (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id) AS likes_count,
                   (SELECT COUNT(*) FROM message_likes WHERE message_id = m.id AND user_id = ?) AS has_liked,
                   (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id) AS dislikes_count,
                   (SELECT COUNT(*) FROM message_dislikes WHERE message_id = m.id AND user_id = ?) AS has_disliked
            FROM messages m 
            JOIN users u ON m.author_id = u.id 
            WHERE m.topic_id = ? 
            ORDER BY m.created_at DESC`
    }

    rows, err := database.DB.Query(commentsQuery, currentUserID, currentUserID, topicID)
    
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
        "IsOwner":       isOwner,
        "User":          user,
        "Sort":          sortParam,
    }
    RenderTemplate(w, r, "view_topic.html", data)
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
    user, _ := GetLoggedUser(r)
    topicID := r.URL.Query().Get("id")

    var authorID int
    err := database.DB.QueryRow("SELECT author_id FROM topics WHERE id = ?", topicID).Scan(&authorID)

    if err != nil || (user.ID != authorID && user.Role != "admin") {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    database.DB.Exec("DELETE FROM topics WHERE id = ?", topicID)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
    user, _ := GetLoggedUser(r)
    if user.ID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    messageID := r.URL.Query().Get("id")
    topicID := r.URL.Query().Get("topic_id")

    var messageAuthorID int
    err := database.DB.QueryRow("SELECT author_id FROM messages WHERE id = ?", messageID).Scan(&messageAuthorID)
    if err != nil {
        http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
        return
    }

    if user.ID != messageAuthorID && user.Role != "admin" {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    database.DB.Exec("DELETE FROM messages WHERE id = ?", messageID)
    http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
}

func EditMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    user, _ := GetLoggedUser(r)
    if user.ID == 0 {
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

    var messageAuthorID int
    err := database.DB.QueryRow("SELECT author_id FROM messages WHERE id = ?", messageID).Scan(&messageAuthorID)
    if err != nil {
        http.Error(w, "Message introuvable", http.StatusNotFound)
        return
    }

    if user.ID != messageAuthorID && user.Role != "admin" {
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

func UpdateTopicStatusHandler(w http.ResponseWriter, r *http.Request) {
    user, _ := GetLoggedUser(r)
    topicID := r.URL.Query().Get("id")
    newStatus := r.URL.Query().Get("status")

    var authorID int
    err := database.DB.QueryRow("SELECT author_id FROM topics WHERE id = ?", topicID).Scan(&authorID)

    if err != nil || (user.ID != authorID && user.Role != "admin") {
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

    // Vérifier si l'utilisateur a déjà épinglé ce topic
    var exists bool
    err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_pins WHERE user_id = ? AND topic_id = ?)", userID, topicID).Scan(&exists)
    if err != nil {
        http.Error(w, "Erreur BDD", 500)
        return
    }

    if exists {
        // Si déjà épinglé, on le retire des favoris
        _, err = database.DB.Exec("DELETE FROM user_pins WHERE user_id = ? AND topic_id = ?", userID, topicID)
    } else {
        // Sinon, on l'ajoute
        _, err = database.DB.Exec("INSERT INTO user_pins (user_id, topic_id) VALUES (?, ?)", userID, topicID)
    }

    if err != nil {
        http.Error(w, "Erreur lors de la modification du favori", 500)
        return
    }
    
    // Rediriger vers la page précédente ou l'accueil
    http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

type EditTopicData struct {
    Topic         models.Topic
    CurrentUserID int
    User          models.User
}

func EditTopicHandler(w http.ResponseWriter, r *http.Request) {
    user, _ := GetLoggedUser(r)
    if user.ID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Récupération de l'ID du topic
    topicID := r.URL.Query().Get("id")
    if topicID == "" {
        topicID = r.FormValue("topic_id")
    }

    // Récupération du topic
    var t models.Topic
    var image sql.NullString

    err := database.DB.QueryRow(`
        SELECT id, title, content, category, author_id, image_url
        FROM topics WHERE id = ?`,
        topicID,
    ).Scan(
        &t.ID, &t.Title, &t.Content, &t.Category, &t.AuthorID, &image,
    )

    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if image.Valid {
        t.ImageURL = image.String
    }

    // Vérification des droits
    if user.ID != t.AuthorID && user.Role != "admin" {
        http.Error(w, "Action non autorisée", http.StatusForbidden)
        return
    }

    // --------------------------
    //  GET → afficher le formulaire
    // --------------------------
    if r.Method == http.MethodGet {
        data := EditTopicData{
            Topic:         t,
            CurrentUserID: user.ID,
            User:          user,
        }

        RenderTemplate(w, r, "edit_topic.html", data)
        return
    }

    // --------------------------
    //  POST → traiter la modification
    // --------------------------
    if r.Method == http.MethodPost {

        title := r.FormValue("title")
        content := r.FormValue("content")
        category := r.FormValue("category")

        if title == "" || content == "" {
            http.Error(w, "Le titre et le contenu ne peuvent pas être vides.", http.StatusBadRequest)
            return
        }

        // --------------------------
        // 1) Suppression d’image
        // --------------------------
        if r.FormValue("delete_image") == "1" {
            _, _ = database.DB.Exec(`
                UPDATE topics SET image_url = NULL WHERE id = ?
            `, topicID)
        }

        // --------------------------
        // 2) Upload d’une nouvelle image
        // --------------------------
        file, header, err := r.FormFile("image")
        if err == nil && header != nil {
            defer file.Close()

            // Générer un nom unique
            filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

            // Sauvegarde dans /uploads/
            out, err := os.Create("uploads/" + filename)
            if err == nil {
                defer out.Close()
                io.Copy(out, file)

                // Mise à jour en BDD
                _, _ = database.DB.Exec(`
                    UPDATE topics SET image_url = ? WHERE id = ?
                `, filename, topicID)
            }
        }

        // --------------------------
        // 3) Mise à jour du texte
        // --------------------------
        _, err = database.DB.Exec(`
            UPDATE topics SET title = ?, content = ?, category = ?
            WHERE id = ?
        `, title, content, category, topicID)

        if err != nil {
            http.Error(w, "Erreur lors de la mise à jour", 500)
            return
        }

        // Redirection finale
        http.Redirect(w, r, "/topic/view?id="+topicID, http.StatusSeeOther)
        return
    }
}
