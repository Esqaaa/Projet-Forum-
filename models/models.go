package models

import "time"

type Topic struct {
	ID 	  		int      	// ID unique du topic
	Title 		string 		// Titre du topic
	Content 	string 		// Contenu du topic
	Tags		string 		// Tags du topic
	CreatedAt 	time.Time 	// Date de création du topic
	AuthorID 	int 		// ID de l'auteur du topic
	Status 		string 		// Statut du topic (ex: "Ouvert", "Fermé", "Archivé")
}

type Message struct {
	ID 	  		int      	// ID unique du message
	TopicID     int         // ID du topic qui contient le message
	AuthorID    int         // ID de l'auteur du message
	Content     string      // Contenu du message
	CreatedAt   time.Time   // Date de création du message
}

