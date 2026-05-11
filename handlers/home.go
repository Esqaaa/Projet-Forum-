package handlers

import (
	"forum/database"
	"net/http"
)

type Topic struct {
	Title   string
	Content string
	Date    string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	_ = cookie

	rows, err := database.DB.Query(
		"SELECT title, content, created_at FROM topics ORDER BY created_at DESC",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var topics []Topic

	for rows.Next() {
		var t Topic
		var date string

		rows.Scan(&t.Title, &t.Content, &date)
		t.Date = date

		topics = append(topics, t)
	}

	RenderTemplate(w, "index.html", topics)
}