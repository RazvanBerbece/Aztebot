package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type TimeoutsRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewTimeoutsRepository() *TimeoutsRepository {
	repo := new(TimeoutsRepository)
	repo.Conn.Connect()
	return repo
}

func (r TimeoutsRepository) GetUserTimeout(userId string) (*dax.Timeout, error) {

	query := "SELECT * FROM Timeouts WHERE userId = ?"
	row := r.Conn.SqlDb.QueryRow(query, userId)

	var timeout dax.Timeout
	err := row.Scan(&timeout.Id,
		&timeout.UserId,
		&timeout.Reason,
		&timeout.CreationTimestamp,
		&timeout.SDuration,
	)

	if err != nil {
		return nil, err
	}

	return &timeout, nil

}

func (r TimeoutsRepository) SaveTimeout(userId string, reason string, timestamp int64, sDuration int) error {

	warn := &dax.Timeout{
		UserId:            userId,
		Reason:            reason,
		CreationTimestamp: timestamp,
		SDuration:         sDuration,
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
		INSERT INTO 
			Timeouts(
				userId, 
				reason, 
				creationTimestamp,
				sTimeLength
			)
		VALUES(?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(warn.UserId, warn.Reason, warn.CreationTimestamp, warn.SDuration)
	if err != nil {
		return err
	}

	return nil

}

func (r TimeoutsRepository) GetTimeoutsCountForUser(userId string) int {
	query := "SELECT COUNT(*) FROM Timeouts WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r TimeoutsRepository) GetArchivedTimeoutsCountForUser(userId string) int {
	query := "SELECT COUNT(*) FROM TimeoutsArchive WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r TimeoutsRepository) ClearTimeoutForUser(userId string) error {

	query := "DELETE FROM Timeouts WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user timeout: %w", err)
	}

	return nil
}

func (r TimeoutsRepository) ClearArchivedTimeout(archivedTimeoutId int64) error {

	query := "DELETE FROM TimeoutsArchive WHERE id = ?"

	_, err := r.Conn.SqlDb.Exec(query, archivedTimeoutId)
	if err != nil {
		return fmt.Errorf("error deleting archived user timeout: %w", err)
	}

	return nil
}

func (r TimeoutsRepository) ClearArchivedTimeoutsForUser(userId string) error {

	query := "DELETE FROM TimeoutsArchive WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting archived users' timeouts: %w", err)
	}

	return nil
}

func (r TimeoutsRepository) GetAllTimeouts() ([]dax.Timeout, error) {

	var timeouts []dax.Timeout

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM Timeouts ORDER BY creationTimestamp ASC")
	if err != nil {
		return nil, fmt.Errorf("GetAllTimeouts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var timeout dax.Timeout
		if err := rows.Scan(&timeout.Id, &timeout.UserId, &timeout.Reason, &timeout.CreationTimestamp, &timeout.SDuration); err != nil {
			return nil, fmt.Errorf("GetAllTimeouts: %v", err)
		}
		timeouts = append(timeouts, timeout)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllTimeouts: %v", err)
	}

	return timeouts, nil

}

func (r TimeoutsRepository) ArchiveTimeout(userId string, reason string, expiryTimestamp int64) error {

	expiredTimeout := &dax.ArchivedTimeout{
		UserId:          userId,
		Reason:          reason,
		ExpiryTimestamp: expiryTimestamp,
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
		INSERT INTO 
			TimeoutsArchive(
				userId, 
				reason, 
				expiryDate
			)
		VALUES(?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(expiredTimeout.UserId, expiredTimeout.Reason, expiredTimeout.ExpiryTimestamp)
	if err != nil {
		return err
	}

	return nil

}

func (r TimeoutsRepository) GetAllArchivedTimeouts() ([]dax.ArchivedTimeout, error) {

	var timeouts []dax.ArchivedTimeout

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM TimeoutsArchive ORDER BY expiryDate ASC")
	if err != nil {
		return nil, fmt.Errorf("GetAllArchivedTimeouts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var timeout dax.ArchivedTimeout
		if err := rows.Scan(&timeout.Id, &timeout.UserId, &timeout.Reason, &timeout.ExpiryTimestamp); err != nil {
			return nil, fmt.Errorf("GetAllArchivedTimeouts: %v", err)
		}
		timeouts = append(timeouts, timeout)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllArchivedTimeouts: %v", err)
	}

	return timeouts, nil

}

func (r TimeoutsRepository) GetAllArchivedTimeoutsForUser(userId string) ([]dax.ArchivedTimeout, error) {

	var timeouts []dax.ArchivedTimeout

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM TimeoutsArchive WHERE userId = ? ORDER BY expiryDate ASC", userId)
	if err != nil {
		return nil, fmt.Errorf("GetAllArchivedTimeoutsForUser: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var timeout dax.ArchivedTimeout
		if err := rows.Scan(&timeout.Id, &timeout.UserId, &timeout.Reason, &timeout.ExpiryTimestamp); err != nil {
			return nil, fmt.Errorf("GetAllArchivedTimeoutsForUser: %v", err)
		}
		timeouts = append(timeouts, timeout)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllArchivedTimeoutsForUser: %v", err)
	}

	return timeouts, nil

}
func (r TimeoutsRepository) UpdateAdminTimeoutCount(adminId string) error {
	_, err := r.Conn.SqlDb.Exec("INSERT INTO ActivitateAdmin (admin_id, timeout_count) VALUES (?, 1) ON DUPLICATE KEY UPDATE timeout_count = timeout_count + 1", adminId)
	return err
}
