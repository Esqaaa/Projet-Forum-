package database

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error

    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN manquant dans les variables d'environnement")
    }

    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Erreur connexion DB:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("DB inaccessible:", err)
    }

    DB.SetMaxOpenConns(10)
    DB.SetMaxIdleConns(5)

    log.Println("[!] Connecté à PostgreSQL")
}
