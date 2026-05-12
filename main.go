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
	http.HandleFunc("/topic/view", handlers.ViewTopicHandler) // Ajout de la route manquante

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Requête SQL sans image_url, avec jointure pour le pseudo
		query := `
			SELECT t.id, t.title, t.content, t.status, t.created_at, u.username 
			FROM topics t 
			JOIN users u ON t.author_id = u.id 
			ORDER BY t.created_at DESC`
		
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
			
			// Scan dans l'ordre du SELECT
			err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Status, &rawDate, &t.Author)
			if err != nil {
				continue 
			}
			
			t.CreatedAt = string(rawDate) // On remplit CreatedAt (ton champ string)
			topics = append(topics, t)
		}

		handlers.RenderTemplate(w, "index.html", topics)
	})

	fmt.Println("[!] Serveur démarré sur http://localhost:8080/register")
	http.ListenAndServe(":8080", nil)
}