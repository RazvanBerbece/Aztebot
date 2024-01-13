package repositories

import (
	"fmt"

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

func (r WarnsRepository) GetWarnWithIdForUser(warnId int64, userId string) (*dataModels.Warn, error) {

	query := "SELECT * FROM Warns WHERE id = ? AND userId = ?"
	row := r.Conn.Db.QueryRow(query, warnId, userId)

	var warn dataModels.Warn
	err := row.Scan(&warn.Id,
		&warn.UserId,
		&warn.Reason,
		&warn.CreationTimestamp,
	)

	if err != nil {
		return nil, err
	}

	return &warn, nil

}

func (r WarnsRepository) GetOldestWarnForUser(userId string) (*dataModels.Warn, error) {

	query := "SELECT * FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC LIMIT 1"
	row := r.Conn.Db.QueryRow(query, userId)

	var warn dataModels.Warn
	err := row.Scan(&warn.Id,
		&warn.UserId,
		&warn.Reason,
		&warn.CreationTimestamp,
	)

	if err != nil {
		return nil, err
	}

	return &warn, nil

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
				creationTimestamp
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

func (r WarnsRepository) GetWarningsCountForUser(userId string) int {
	query := "SELECT COUNT(*) FROM Warns WHERE userId = ?"
	var count int
	err := r.Conn.Db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r WarnsRepository) DeleteAllWarningsForUser(userId string) error {

	query := "DELETE FROM Warns WHERE userId = ?"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting all user warnings: %w", err)
	}

	return nil
}

func (r WarnsRepository) DeleteOldestWarningForUser(userId string) error {

	query := "DELETE FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC LIMIT 1"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting oldest user warning: %w", err)
	}

	return nil
}

func (r WarnsRepository) GetWarningsForUser(userId string) ([]dataModels.Warn, error) {

	var warns []dataModels.Warn

	rows, err := r.Conn.Db.Query("SELECT * FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC", userId)
	if err != nil {
		return nil, fmt.Errorf("GetWarningsForUser %s: %v", userId, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var warn dataModels.Warn
		if err := rows.Scan(&warn.Id, &warn.UserId, &warn.Reason, &warn.CreationTimestamp); err != nil {
			return nil, fmt.Errorf("GetWarningsForUser %s: %v", userId, err)
		}
		warns = append(warns, warn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetWarningsForUser %s: %v", userId, err)
	}

	// Check for zero rows
	if len(warns) == 0 {
		return []dataModels.Warn{}, nil
	}

	return warns, nil
}
