package handlers

import (
    "forum/database"
    "forum/models"
    "net/http"
)

// Les accès admin sur un dashboard séparé  
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
    user, err := GetLoggedUser(r)
    if err != nil || user.Role != "admin" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    rows, _ := database.DB.Query(`
        SELECT id, title, author_id
        FROM topics
        ORDER BY created_at DESC
    `)

    var topics []models.Topic
    for rows.Next() {
        var t models.Topic
        rows.Scan(&t.ID, &t.Title, &t.AuthorID)
        topics = append(topics, t)
    }

    userRows, _ := database.DB.Query(`
        SELECT id, username, email, role
        FROM users
        ORDER BY id ASC
    `)

    var users []models.User
    for userRows.Next() {
        var u models.User
        userRows.Scan(&u.ID, &u.Username, &u.Email, &u.Role)
        users = append(users, u)
    }

    RenderTemplate(w, r, "admin_dashboard.html", map[string]interface{}{
        "Topics": topics,
        "Users":  users,
    })
}

// La fonction de bannissement des users 
func AdminBanUser(w http.ResponseWriter, r *http.Request) {
    user, err := GetLoggedUser(r)
    if err != nil || user.Role != "admin" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    username := r.URL.Query().Get("username")

    database.DB.Exec("DELETE FROM users WHERE username = ?", username)

    http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
