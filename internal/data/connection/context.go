package databaseconn

import (
	"database/sql"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseContext interface {
	Connect()
}

type AztebotDbContext struct {
	ConnectionString string
	SqlDb            *sql.DB
}

type AztemarketDbContext struct {
	ConnectionString string
	SqlDb            *sql.DB
}

func (c *AztebotDbContext) Connect() {

	db, err := sql.Open("mysql", globalConfiguration.MySqlAztebotRootConnectionString)
	if err != nil {
		log.Fatal("Connection to AzteBot database cannot be established :", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database at %s cannot be reached : %s", globalConfiguration.MySqlAztebotRootConnectionString, err)
	}

	c.SqlDb = db

}

func (c *AztemarketDbContext) Connect() {

	db, err := sql.Open("mysql", globalConfiguration.MySqlAztemarketRootConnectionString)
	if err != nil {
		log.Fatal("Connection to AzteMarket database cannot be established :", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database at %s cannot be reached : %s", globalConfiguration.MySqlAztemarketRootConnectionString, err)
	}

	c.SqlDb = db

}
