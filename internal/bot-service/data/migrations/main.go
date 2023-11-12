package main

import (
	"database/sql"
	"embed"
	"log"

	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/globals"
	"github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed history/*.sql
var embedMigrations embed.FS

func main() {

	db, err := sql.Open("mysql", globals.MySqlRootPrivateTcpConnectionString)
	if err != nil {
		log.Fatal("Connection to database cannot be established:", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("Database cannot be reached:", pingErr)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatal("Cannot set mysql dialect for Goose DB migration:", err)
	}

	if err := goose.Up(db, "history"); err != nil {
		log.Fatal("Cannot run UP migration:", err)
	}
}
