package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/models" // Import des models indispensable
	"net/http"
)

func main() {
	database.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/topic/create", handlers.CreateTopicHandler)
	http.HandleFunc("/topic/view", handlers.ViewTopicHandler)
	http.HandleFunc("/topic/pin", handlers.PinTopicHandler)
	http.HandleFunc("/message/post", handlers.PostMessageHandler)
	http.HandleFunc("/topic/delete", handlers.DeleteTopicHandler)
	http.HandleFunc("/message/delete", handlers.DeleteMessageHandler)
	http.HandleFunc("/topic/update-status", handlers.UpdateTopicStatusHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Si l'URL n'est pas exactement "/", on renvoie une erreur 404
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Protection : Si l'utilisateur n'est pas connecté, on le dégage vers le login
		userID := handlers.GetLoggedUserID(r) 
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Requête SQL sans image_url, avec jointure pour le pseudo
		query := `
			SELECT t.id, t.title, t.content, t.status, t.is_pinned, t.created_at, u.username 
			FROM topics t 
			JOIN users u ON t.author_id = u.id 
			ORDER BY t.is_pinned DESC, t.created_at DESC`
		
		rows, err := database.DB.Query(query)
		if err != nil {
			http.Error(w, "Erreur BDD : "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var topics []models.Topic 
		for rows.Next() {
			var t models.Topic
			var rawDate []byte
			
			err := rows.Scan(
				&t.ID, 
				&t.Title, 
				&t.Content, 
				&t.Status, 
				&t.IsPinned, 
				&rawDate, 
				&t.Author,
			)
			if err != nil {
				fmt.Println("Erreur scan topic:", err)
				continue 
			}
			
			t.CreatedAt = string(rawDate)
			topics = append(topics, t)
		}

		handlers.RenderTemplate(w, "index.html", topics)
	})

	fmt.Println("[!] Serveur démarré sur http://localhost:8080/register")
	http.ListenAndServe(":8080", nil)
}