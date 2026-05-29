package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"net/http"
)

func main() {
	database.InitDB()

	// Gestion des fichiers statiques (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes de l'application
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/topic/create", handlers.CreateTopicHandler)
	http.HandleFunc("/topic/view", handlers.ViewTopicHandler)
	http.HandleFunc("/topic/pin", handlers.PinTopicHandler)
	http.HandleFunc("/message/post", handlers.PostMessageHandler)
	http.HandleFunc("/topic/delete", handlers.DeleteTopicHandler)
	http.HandleFunc("/message/edit", handlers.EditMessageHandler)
	http.HandleFunc("/message/delete", handlers.DeleteMessageHandler)
	http.HandleFunc("/topic/update-status", handlers.UpdateTopicStatusHandler)
	http.HandleFunc("/message/like", handlers.LikeMessageHandler)

	http.HandleFunc("/", handlers.HomeHandler)

	fmt.Println("[!] Serveur démarré sur http://localhost:8080/register")
	http.ListenAndServe(":8080", nil)
}