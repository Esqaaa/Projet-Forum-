package models

type User struct {
	ID 			int 		// Identifiant unique 
	Username 	string 		// Pseudo de l'utilisateur connecté 
	Email 		string 		// Email de l'utilisateur connecté 
	Password 	string 		// Mot de passe de l'utilisateur connecté
	Role		string 		// 'user' ou 'admin' 
}

type UserProfile struct {
	ID           int    	// ID unique de l'utilisateur
	Username     string 	// Nom d'utilisateur
	Email        string 	// Adresse e-mail de l'utilisateur
	Biography    string 	// Biographie de l'utilisateur
	AvatarURL    string 	// URL de l'avatar de l'utilisateur
	LastLogin    string 	// Date de la dernière connexion de l'utilisateur
	TopicsCount  int    	// Nombre total de topics créés par l'utilisateur
	CommentCount int    	// Nombre total de commentaires postés par l'utilisateur
}
