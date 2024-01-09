package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
)

type WarnsRepository struct {
	Conn databaseconn.Database
}

func NewWarnsRepository() *WarnsRepository {
	repo := new(WarnsRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}
