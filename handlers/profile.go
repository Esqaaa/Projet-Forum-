package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"io"
	"net/http"
	"os"
	"time"
)

// ProfileHandler affiche le profil et ses paramètres
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID := GetLoggedUserID(r)
	if currentUserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var p models.UserProfile
	var rawDate []byte
	var bio sql.NullString
	var avatar sql.NullString

	// Récupération des infos de l'utilisateur
	queryUser := "SELECT id, username, email, biography, avatar_url, last_login FROM users WHERE id = ?"
	err := database.DB.QueryRow(queryUser, currentUserID).Scan(&p.ID, &p.Username, &p.Email, &bio, &avatar, &rawDate)
	if err != nil {
		fmt.Println("Erreur récupération profil :", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Gestion des valeurs NULL ou vides pour l'affichage
	if bio.Valid { p.Biography = bio.String }
	if avatar.Valid && avatar.String != "" { 
		p.AvatarURL = avatar.String 
	} else {
		p.AvatarURL = "/static/uploads/default-avatar.png"
	}
	p.LastLogin = string(rawDate)

	// Calcul du nombre de topics créés
	database.DB.QueryRow("SELECT COUNT(*) FROM topics WHERE author_id = ?", currentUserID).Scan(&p.TopicsCount)

	// Calcul du nombre de messages envoyés
	database.DB.QueryRow("SELECT COUNT(*) FROM messages WHERE author_id = ?", currentUserID).Scan(&p.CommentCount)

	status := r.URL.Query().Get("status")

	data := map[string]interface{}{
		"Profile":       p,
		"CurrentUserID": currentUserID,
		"Status":        status,
	}

	RenderTemplate(w, r, "profile.html", data)
}

// UpdateProfileHandler enregistre les modifications (Bio, Email, Photo)
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	currentUserID := GetLoggedUserID(r)
	if currentUserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// On autorise l'upload d'images jusqu'à 3 Mo
	r.ParseMultipartForm(3 << 20)

	newEmail := r.FormValue("email")
	newBio := r.FormValue("biography")

	// Récupération de l'ancienne image pour ne pas l'écraser si l'utilisateur n'envoie rien
	var currentAvatar string
	database.DB.QueryRow("SELECT avatar_url FROM users WHERE id = ?", currentUserID).Scan(&currentAvatar)

	// Gestion de l'upload de la photo de profil
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		
		// On crée un nom unique pour l'image
		fileName := fmt.Sprintf("avatar-%d-%s", time.Now().Unix(), handler.Filename)
		imagePath := "/static/uploads/" + fileName
		
		// On s'assure que le dossier static/uploads existe
		os.MkdirAll("./static/uploads", os.ModePerm)

		dst, err := os.Create("./" + imagePath)
		if err == nil {
			defer dst.Close()
			io.Copy(dst, file)
			currentAvatar = imagePath // On remplace par le nouveau chemin
		}
	}

	// Mise à jour dans la base de données
	query := "UPDATE users SET email = ?, biography = ?, avatar_url = ? WHERE id = ?"
	_, err = database.DB.Exec(query, newEmail, newBio, currentAvatar, currentUserID)
	if err != nil {
		fmt.Println("Erreur update profil :", err)
		http.Redirect(w, r, "/profile?status=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile?status=success", http.StatusSeeOther)
}