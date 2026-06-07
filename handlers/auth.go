package handlers

import (
    "database/sql"
    "errors"
    "forum/database"
    "net/http"
    "regexp"
    "html/template"

    "golang.org/x/crypto/bcrypt"
)

type TemplateData struct {
    Errors     []string
    Identifier string
    Username   string
    Email      string
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
    // Récupérer l'utilisateur connecté
    user, _ := GetLoggedUser(r)

    if data == nil {
        data = map[string]interface{}{}
    }

    data["User"] = user

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
        RenderTemplate(w, r, "register.html", nil)
        return
    }

    r.ParseForm()

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    var errorsList []string

    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
    if !usernameRegex.MatchString(username) {
        errorsList = append(errorsList, "Pseudo invalide")
    }

    emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
    if !emailRegex.MatchString(email) {
        errorsList = append(errorsList, "Email invalide")
    }

    if err := ValidatePassword(password); err != nil {
        errorsList = append(errorsList, err.Error())
    }

    if len(errorsList) > 0 {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   errorsList,
            "Username": username,
            "Email":    email,
        })
        return
    }

    var tmp int
    err := database.DB.QueryRow(
        "SELECT id FROM users WHERE username = ?",
        username,
    ).Scan(&tmp)

    if err == nil {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Pseudo déjà utilisé"},
            "Username": username,
            "Email":    email,
        })
        return
    }

    if err != sql.ErrNoRows {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Erreur serveur DB"},
            "Username": username,
            "Email":    email,
        })
        return
    }

    err = database.DB.QueryRow(
        "SELECT id FROM users WHERE email = ?",
        email,
    ).Scan(&tmp)

    if err == nil {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Email déjà utilisé"},
            "Username": username,
            "Email":    email,
        })
        return
    }

    if err != sql.ErrNoRows {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Erreur serveur DB"},
            "Username": username,
            "Email":    email,
        })
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Erreur hash password"},
            "Username": username,
            "Email":    email,
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
        RenderTemplate(w, r, "register.html", map[string]interface{}{
            "Errors":   []string{"Erreur insertion utilisateur"},
            "Username": username,
            "Email":    email,
        })
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        RenderTemplate(w, r, "login.html", nil)
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
        RenderTemplate(w, r, "login.html", map[string]interface{}{
            "Errors":     []string{"Identifiants invalides"},
            "Identifier": identifier,
        })
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
    if err != nil {
        RenderTemplate(w, r, "login.html", map[string]interface{}{
            "Errors":     []string{"Mot de passe incorrect"},
            "Identifier": identifier,
        })
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
