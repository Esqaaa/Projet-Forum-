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

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(
		"templates/layout.html",
		"templates/"+tmpl,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "layout", data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		RenderTemplate(w, "register.html", nil)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validation pseudo
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !usernameRegex.MatchString(username) {
		http.Error(w, "Pseudo invalide", http.StatusBadRequest)
		return
	}

	// Validation password
	err := ValidatePassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifie username unique
	var exists int

	err = database.DB.QueryRow(
		"SELECT id FROM users WHERE username = ?",
		username,
	).Scan(&exists)

	if err != sql.ErrNoRows {
		http.Error(w, "Pseudo déjà utilisé", http.StatusBadRequest)
		return
	}

	// Vérifie email unique
	err = database.DB.QueryRow(
		"SELECT id FROM users WHERE email = ?",
		email,
	).Scan(&exists)

	if err != sql.ErrNoRows {
		http.Error(w, "Email déjà utilisé", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Insert user
	_, err = database.DB.Exec(
		"INSERT INTO users(username, email, password) VALUES (?, ?, ?)",
		username,
		email,
		string(hashedPassword),
	)

	if err != nil {
		http.Error(w, "Erreur insertion", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		RenderTemplate(w, "login.html", nil)
		return
	}

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	var id int
	var hashedPassword string

	err := database.DB.QueryRow(
		`SELECT id, password
		FROM users
		WHERE username = ? OR email = ?`,
		identifier,
		identifier,
	).Scan(&id, &hashedPassword)

	if err != nil {
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	if err != nil {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Cookie session simple
	cookie := &http.Cookie{
		Name:  "session",
		Value: identifier,
		Path:  "/",
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ValidatePassword(password string) error {

	if len(password) < 8 {
		return errors.New("Le mot de passe doit faire 8 caractères minimum")
	}

	uppercase := regexp.MustCompile(`[A-Z]`)
	if !uppercase.MatchString(password) {
		return errors.New("Le mot de passe doit contenir une majuscule")
	}

	special := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	if !special.MatchString(password) {
		return errors.New("Le mot de passe doit contenir un caractère spécial")
	}

	return nil
}