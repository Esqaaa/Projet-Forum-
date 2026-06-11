package handlers

import (
    "database/sql"
    "errors"
    "forum/database"
    "html/template"
    "net/http"
    "regexp"
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

type TemplateData struct {
    Errors     []string
    Identifier string
    Username   string
    Email      string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
    t, err := template.ParseFiles(
        "templates/html/layout.html",
        "templates/html/"+tmpl,
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

    var errorsList []string

    // VALIDATION USERNAME
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
    if !usernameRegex.MatchString(username) {
        errorsList = append(errorsList, "Pseudo invalide")
    }

    // VALIDATION EMAIL
    emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
    if !emailRegex.MatchString(email) {
        errorsList = append(errorsList, "Email invalide")
    }

    // VALIDATION PASSWORD
    if err := ValidatePassword(password); err != nil {
        errorsList = append(errorsList, err.Error())
    }

    // SI ERREURS → AFFICHER
    if len(errorsList) > 0 {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   errorsList,
            Username: username,
            Email:    email,
        })
        return
    }

    // VERIFICATION USERNAME EXISTE
    var tmp int
    err := database.DB.QueryRow(
        "SELECT id FROM users WHERE BINARY username = ?",
        username,
    ).Scan(&tmp)

    if err == nil {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Pseudo déjà utilisé"},
            Username: username,
            Email:    email,
        })
        return
    }

    if err != sql.ErrNoRows {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Erreur serveur DB"},
            Username: username,
            Email:    email,
        })
        return
    }

    // VERIFICATION EMAIL EXISTE
    err = database.DB.QueryRow(
        "SELECT id FROM users WHERE BINARY email = ?",
        email,
    ).Scan(&tmp)

    if err == nil {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Email déjà utilisé"},
            Username: username,
            Email:    email,
        })
        return
    }

    if err != sql.ErrNoRows {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Erreur serveur DB"},
            Username: username,
            Email:    email,
        })
        return
    }

    // HASH PASSWORD
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Erreur hash password"},
            Username: username,
            Email:    email,
        })
        return
    }

    _, err = database.DB.Exec(
        "INSERT INTO users(username, email, password) VALUES (?, ?, ?)",
        username,
        email,
        string(hashedPassword),
    )

    if err != nil {
        RenderTemplate(w, "register.html", TemplateData{
            Errors:   []string{"Erreur insertion utilisateur"},
            Username: username,
            Email:    email,
        })
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
    var officialUsername string

    err := database.DB.QueryRow(
        `SELECT id, password, username FROM users WHERE BINARY username = ? OR BINARY email = ?`,
        identifier,
        identifier,
    ).Scan(&id, &hashed, &officialUsername)

    if err != nil {
        RenderTemplate(w, "login.html", TemplateData{
            Errors:     []string{"Identifiants invalides"},
            Identifier: identifier,
        })
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
    if err != nil {
        RenderTemplate(w, "login.html", TemplateData{
            Errors:     []string{"Mot de passe incorrect"},
            Identifier: identifier,
        })
        return
    }

    // Mise à jour du last_login
    _, err = database.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?", id)
    if err != nil {
        fmt.Println("Erreur lors de la mise à jour de last_login :", err)
    }

    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: officialUsername,
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
