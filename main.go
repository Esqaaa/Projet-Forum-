package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"net/http"
)

func main() {

	database.InitDB()

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/topic/create", handlers.CreateTopicHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. On va chercher les topics en base de données
		rows, err := database.DB.Query("SELECT title, content, created_at FROM topics ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "Erreur BDD : "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 2. On définit la structure locale pour stocker les données
		type Topic struct {
			Title   string
			Content string
			Date    string
		}
		var topics []Topic

		// 3. On remplit la liste "topics"
		for rows.Next() {
			var t Topic
			var rawDate []byte // Le driver SQL renvoie souvent la date en bytes
			rows.Scan(&t.Title, &t.Content, &rawDate)
			t.Date = string(rawDate)
			topics = append(topics, t)
		}

		// 4. ON ENVOIE ENFIN LA VARIABLE 'topics' AU LIEU DE 'nil'
		handlers.RenderTemplate(w, "index.html", topics)
	})

	fmt.Println("Serveur lancé sur http://localhost:8080/register")
	http.ListenAndServe(":8080", nil)
}