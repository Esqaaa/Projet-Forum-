package models

import "time"

type Topic struct {
	ID 	  		int      	// ID unique du topic
	Title 		string 		// Titre du topic
	Content 	string 		// Contenu du topic
	Tags		string 		// Tags du topic
	CreatedAt 	string 	    // Date de création du topic
	Date        string      // Date formatée pour l'affichage HTML
	AuthorID 	int 		// ID de l'auteur du topic
	Author      string      // Nom de l'auteur du topic
	Status 		string 		// Statut du topic (ex: "Ouvert", "Fermé", "Archivé")
}

type Message struct {
	ID 	  		int      	// ID unique du message
	TopicID     int         // ID du topic qui contient le message
	AuthorID    int         // ID de l'auteur du message
	Content     string      // Contenu du message
	CreatedAt   time.Time   // Date de création du message
}

type Comment struct {
    ID        int
    Content   string
    TopicID   int
    Author    string 
    Date      string
}

