package main

import (
	"embed"

	databasePackage "github.com/LxrdVixxeN/Aztebot/internal/bot-service/data/connection"
	"github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed history/*.sql
var embedMigrations embed.FS

func main() {

	var database databasePackage.Database
	database.ConnectDatabaseHandle()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	if err := goose.Up(database.Db, "history"); err != nil {
		panic(err)
	}
}
