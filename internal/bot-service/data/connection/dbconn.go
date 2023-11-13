package databaseconn

import (
	"database/sql"
	"log"

	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/globals"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Db *sql.DB
}

func (d *Database) ConnectDatabaseHandle() {

	db, err := sql.Open("mysql", globals.MySqlRootConnectionString)
	if err != nil {
		log.Fatal("Connection to database cannot be established :", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("Database at %s cannot be reached : %s", globals.MySqlRootConnectionString, pingErr)
	}

	d.Db = db

}
