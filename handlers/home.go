package handlers

import (
    "database/sql"
    "fmt"
    "forum/database"
    "forum/models"
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    currentUserID := GetLoggedUserID(r)
    if currentUserID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    pageSize := 4
    currentPage := 1

    pageStr := r.URL.Query().Get("page")
    if pageStr != "" {
        fmt.Sscanf(pageStr, "%d", &currentPage)
    }

    if currentPage < 1 {
        currentPage = 1
    }

    categoryFilter := r.URL.Query().Get("category")

    var totalTopics int
    var err error 
    if categoryFilter != "" {
        err = database.DB.QueryRow("SELECT COUNT(*) FROM topics WHERE category = ?", categoryFilter).Scan(&totalTopics)
    } else {
        err = database.DB.QueryRow("SELECT COUNT(*) FROM topics").Scan(&totalTopics)
    }
    if err != nil {
        totalTopics = 0
    }

    totalPages := (totalTopics + pageSize - 1) / pageSize
    if totalPages < 1 {
        totalPages = 1
    }
    if currentPage > totalPages {
        currentPage = totalPages
    }

    offset := (currentPage - 1) * pageSize

    queryStr := `
        SELECT 
            t.id, t.title, t.content, t.created_at, t.status, 
            t.is_pinned, t.image_url, t.category, u.username, u.id
        FROM topics t
        JOIN users u ON t.author_id = u.id `

    var rows *sql.Rows
    if categoryFilter != "" {
        queryStr += `WHERE t.category = ? 
                     ORDER BY t.is_pinned DESC, t.created_at DESC 
                     LIMIT ? OFFSET ?`
        rows, err = database.DB.Query(queryStr, categoryFilter, pageSize, offset)
    } else {
        queryStr += `ORDER BY t.is_pinned DESC, t.created_at DESC 
                     LIMIT ? OFFSET ?`
        rows, err = database.DB.Query(queryStr, pageSize, offset)
    }

    if err != nil {
        http.Error(w, "Erreur récupération topics : "+err.Error(), 500)
        return
    }
    defer rows.Close()

    var topics []models.Topic
    for rows.Next() {
        var t models.Topic
        var rawDate []byte
        var imageURL sql.NullString

        err := rows.Scan(
            &t.ID, &t.Title, &t.Content, &rawDate, &t.Status,
            &t.IsPinned, &imageURL, &t.Category, &t.Author, &t.AuthorID,
        )
        if err != nil {
            fmt.Println("Erreur scan topic:", err)
            continue
        }

        if imageURL.Valid {
            t.ImageURL = imageURL.String
        } else {
            t.ImageURL = ""
        }

        t.Date = string(rawDate)
        t.CreatedAt = string(rawDate)
        topics = append(topics, t)
    }

    categories := []string{"Sport", "Musique", "Automobile", "Aviation", "Sciences", "Informatique"}
    
    data := map[string]interface{}{
        "Topics":        topics,
        "CurrentUserID": currentUserID,
        "CurrentPage": currentPage,
        "TotalPages":  totalPages,
        "HasPrev":     currentPage > 1,
        "HasNext":     currentPage < totalPages,
        "PrevPage":    currentPage - 1,
        "NextPage":    currentPage + 1,
        "SelectedCategory": categoryFilter,
        "Categories": categories,
    }

    RenderTemplate(w, "index.html", data)
}