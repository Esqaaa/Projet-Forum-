package models

type User struct {
	ID 			int 		// Identifiant unique 
	Username 	string 		// Pseudo de l'utilisateur connecté 
	Email 		string 		// Email de l'utilisateur connecté 
	Password 	string 		// Mot de passe de l'utilisateur connecté 
}