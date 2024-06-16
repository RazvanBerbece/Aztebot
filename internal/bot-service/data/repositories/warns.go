package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type WarnsRepository struct {
	Conn databaseconn.Database
}

func NewWarnsRepository() *WarnsRepository {
	repo := new(WarnsRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r WarnsRepository) SaveWarn(userId string, reason string, timestamp int64) error {

	warn := &dataModels.Warn{
		UserId:            userId,
		Reason:            reason,
		CreationTimestamp: timestamp,
	}

	stmt, err := r.Conn.Db.Prepare(`
		INSERT INTO 
			Warns(
				userId, 
				reason, 
				creationTimestamp,
			)
		VALUES(?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(warn.UserId, warn.Reason, warn.CreationTimestamp)
	if err != nil {
		return err
	}

	return nil

}
