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

	// --- AUTHENTIFICATION ---
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// --- TOPICS ---
	http.HandleFunc("/topic/create", handlers.CreateTopicHandler)
	http.HandleFunc("/topic/view", handlers.ViewTopicHandler)
	http.HandleFunc("/topic/edit", handlers.EditTopicHandler)
	http.HandleFunc("/topic/delete", handlers.DeleteTopicHandler)
	http.HandleFunc("/topic/pin", handlers.PinTopicHandler)
	http.HandleFunc("/topic/update-status", handlers.UpdateTopicStatusHandler)

	// --- MESSAGES ---
	http.HandleFunc("/message/post", handlers.PostMessageHandler)
	http.HandleFunc("/message/edit", handlers.EditMessageHandler)
	http.HandleFunc("/message/delete", handlers.DeleteMessageHandler)
	http.HandleFunc("/message/like", handlers.LikeMessageHandler)
	http.HandleFunc("/message/dislike", handlers.DislikeMessageHandler)

	// --- ADMIN ---
	http.HandleFunc("/admin", handlers.AdminDashboard)
	http.HandleFunc("/admin/ban", handlers.AdminBanUser)

	// --- UTILISATEUR ---
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/profile/update", handlers.UpdateProfileHandler)

	// --- RECHERCHE ---
	http.HandleFunc("/search", handlers.SearchHandler)

	// --- HOME ---
	http.HandleFunc("/", handlers.HomeHandler)


	fmt.Println("[!] Serveur démarré sur http://localhost:8080/register")
	http.ListenAndServe(":8080", nil)
}