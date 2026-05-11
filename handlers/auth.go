package handlers

import (
	"database/sql"
	"errors"
	"forum/database"
	"html/template"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles(
		"templates/layout.html",
		"templates/"+tmpl,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		RenderTemplate(w, "register.html", nil)
		return
	}

	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validation username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !usernameRegex.MatchString(username) {
		http.Error(w, "Pseudo invalide", http.StatusBadRequest)
		return
	}

	// Password rules
	if err := ValidatePassword(password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK USERNAME EXIST
	var tmp int
	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE username = ?",
		username,
	).Scan(&tmp)

	if err == nil {
		http.Error(w, "Pseudo déjà utilisé", http.StatusBadRequest)
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "Erreur serveur DB", http.StatusInternalServerError)
		return
	}

	// CHECK EMAIL EXIST
	err = database.DB.QueryRow(
		"SELECT id FROM users WHERE email = ?",
		email,
	).Scan(&tmp)

	if err == nil {
		http.Error(w, "Email déjà utilisé", http.StatusBadRequest)
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "Erreur serveur DB", http.StatusInternalServerError)
		return
	}

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur hash password", http.StatusInternalServerError)
		return
	}

	// INSERT USER
	_, err = database.DB.Exec(
		"INSERT INTO users(username, email, password) VALUES (?, ?, ?)",
		username,
		email,
		string(hashedPassword),
	)

	if err != nil {
		http.Error(w, "Erreur insertion utilisateur", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		RenderTemplate(w, "login.html", nil)
		return
	}

	r.ParseForm()

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	var id int
	var hashed string

	err := database.DB.QueryRow(
		`SELECT id, password FROM users WHERE username = ? OR email = ?`,
		identifier,
		identifier,
	).Scan(&id, &hashed)

	if err != nil {
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: identifier,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ValidatePassword(password string) error {

	if len(password) < 8 {
		return errors.New("8 caractères minimum")
	}

	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return errors.New("1 majuscule requise")
	}

	if ok, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); !ok {
		return errors.New("1 caractère spécial requis")
	}

	return nil
}

