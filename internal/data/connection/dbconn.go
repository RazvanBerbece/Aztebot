package databaseconn

import (
	"database/sql"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Db *sql.DB
}

func (d *Database) ConnectDatabaseHandle() {

	db, err := sql.Open("mysql", globalConfiguration.MySqlRootConnectionString)
	if err != nil {
		log.Fatal("Connection to database cannot be established :", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("Database at %s cannot be reached : %s", globalConfiguration.MySqlRootConnectionString, pingErr)
	}

	d.Db = db

}
