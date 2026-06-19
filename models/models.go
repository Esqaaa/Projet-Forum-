package models

import "time"

type Topic struct {
	ID        int    	// ID unique du topic
	Title     string 	// Titre du topic
	Content   string 	// Contenu du topic
	Tags      string	 // Tags du topic
	CreatedAt string 	// Date de création du topic
	Date      string 	// Date formatée pour l'affichage HTML
	AuthorID  int    	// ID de l'auteur du topic
	Author    string 	// Nom de l'auteur du topic
	IsPinnedByUser bool	// Si le topic est épinglé par le user 
	Category  string 	// Le thème du topic
	ImageURL  string 	// Ajouter une image au post
	Status    string 	// Statut du topic (ex: "Ouvert", "Fermé", "Archivé")
}

type Message struct {
	ID        int       // ID unique du message
	TopicID   int       // ID du topic qui contient le message
	AuthorID  int       // ID de l'auteur du message
	Content   string    // Contenu du message
	CreatedAt time.Time // Date de création du message
}

type Comment struct {
	ID            int    // ID unique du commentaire
	Content       string // Contenu du commentaire
	TopicID       int    // ID du topic auquel le commentaire appartient
	Author        string // Nom de l'auteur du commentaire
	AuthorID      int    // ID de l'auteur
	Date          string // Date d'envoi du commentaire
	LikesCount    int    // Nombre total de like
	HasLiked      bool   // Indique si l'utilisateur connecté a liké
	DislikesCount int    // Nombre total de dislikes
	HasDisliked   bool   // Indique si l'utilisateur connecté a disliké
}
