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

type TemplateData struct {
	Error string
}

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

	// AFFICHAGE PAGE
	if r.Method == "GET" {
		RenderTemplate(w, "register.html", nil)
		return
	}

	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// VALIDATION USERNAME
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if !usernameRegex.MatchString(username) {
		RenderTemplate(w, "register.html", TemplateData{
			Error: "Pseudo invalide",
		})
		return
	}

	// VALIDATION PASSWORD
	if err := ValidatePassword(password); err != nil {

		RenderTemplate(w, "register.html", TemplateData{
			Error: err.Error(),
		})
		return
	}

	// CHECK USERNAME EXIST
	var tmp int

	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE username = ?",
		username,
	).Scan(&tmp)

	if err == nil {

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Pseudo déjà utilisé",
		})
		return
	}

	if err != sql.ErrNoRows {

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Erreur serveur DB",
		})
		return
	}

	// CHECK EMAIL EXIST
	err = database.DB.QueryRow(
		"SELECT id FROM users WHERE email = ?",
		email,
	).Scan(&tmp)

	if err == nil {

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Email déjà utilisé",
		})
		return
	}

	if err != sql.ErrNoRows {

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Erreur serveur DB",
		})
		return
	}

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Erreur hash password",
		})
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

		RenderTemplate(w, "register.html", TemplateData{
			Error: "Erreur insertion utilisateur",
		})
		return
	}

	// REDIRECTION LOGIN
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// AFFICHAGE PAGE
	if r.Method == "GET" {
		RenderTemplate(w, "login.html", nil)
		return
	}

	r.ParseForm()

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	var id int
	var hashed string

	// RECHERCHE USER
	err := database.DB.QueryRow(
		`SELECT id, password FROM users WHERE username = ? OR email = ?`,
		identifier,
		identifier,
	).Scan(&id, &hashed)

	if err != nil {

		RenderTemplate(w, "login.html", TemplateData{
			Error: "Identifiants invalides",
		})
		return
	}

	// CHECK PASSWORD
	err = bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(password),
	)

	if err != nil {

		RenderTemplate(w, "login.html", TemplateData{
			Error: "Mot de passe incorrect",
		})
		return
	}

	// COOKIE SESSION
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: identifier,
		Path:  "/",
	})

	// REDIRECTION ACCUEIL
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