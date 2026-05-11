package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	dsn := "root:@tcp(127.0.0.1:3306)/forum_project"

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur connexion DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB inaccessible:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	log.Println("✔️ Connecté à MySQL")
}