package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type WarnsRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewWarnsRepository() *WarnsRepository {
	repo := new(WarnsRepository)
	repo.Conn.Connect()
	return repo
}

func (r WarnsRepository) GetWarnWithIdForUser(warnId int64, userId string) (*dax.Warn, error) {

	query := "SELECT * FROM Warns WHERE id = ? AND userId = ?"
	row := r.Conn.SqlDb.QueryRow(query, warnId, userId)

	var warn dax.Warn
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

func (r WarnsRepository) GetOldestWarnForUser(userId string) (*dax.Warn, error) {

	query := "SELECT * FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC LIMIT 1"
	row := r.Conn.SqlDb.QueryRow(query, userId)

	var warn dax.Warn
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

func (r WarnsRepository) GetAllWarns() ([]dax.Warn, error) {

	var warns []dax.Warn

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM Warns ORDER BY creationTimestamp ASC")
	if err != nil {
		return nil, fmt.Errorf("GetAllWarns: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var warn dax.Warn
		if err := rows.Scan(&warn.Id, &warn.UserId, &warn.Reason, &warn.CreationTimestamp); err != nil {
			return nil, fmt.Errorf("GetAllWarns: %v", err)
		}
		warns = append(warns, warn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllWarns: %v", err)
	}

	// Check for zero rows
	if len(warns) == 0 {
		return []dax.Warn{}, nil
	}

	return warns, nil

}

func (r WarnsRepository) SaveWarn(userId string, reason string, timestamp int64) error {

	warn := &dax.Warn{
		UserId:            userId,
		Reason:            reason,
		CreationTimestamp: timestamp,
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
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
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r WarnsRepository) DeleteAllWarningsForUser(userId string) error {

	query := "DELETE FROM Warns WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting all user warnings: %w", err)
	}

	return nil
}

func (r WarnsRepository) DeleteOldestWarningForUser(userId string) error {

	query := "DELETE FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC LIMIT 1"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting oldest user warning: %w", err)
	}

	return nil
}

func (r WarnsRepository) GetWarningsForUser(userId string) ([]dax.Warn, error) {

	var warns []dax.Warn

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM Warns WHERE userId = ? ORDER BY creationTimestamp ASC", userId)
	if err != nil {
		return nil, fmt.Errorf("GetWarningsForUser %s: %v", userId, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var warn dax.Warn
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
		return []dax.Warn{}, nil
	}

	return warns, nil
}

func (r WarnsRepository) DeleteWarningForUser(id int64, userId string) error {

	query := "DELETE FROM Warns WHERE userId = ? AND id = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId, id)
	if err != nil {
		return fmt.Errorf("error deleting user warning: %w", err)
	}

	return nil
}
