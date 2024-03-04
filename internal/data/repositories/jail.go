package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
)

type JailRepository struct {
	Conn databaseconn.Database
}

func NewJailRepository() *JailRepository {
	repo := new(JailRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r JailRepository) UserIsJailed(userId string) int {
	query := "SELECT COUNT(*) FROM Jail WHERE userId = ?"
	var count int
	err := r.Conn.Db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for jail entry: %v\n", err)
		return -1
	}
	return count
}

func (r JailRepository) AddUserToJail(userId string, reason string, task string, timestamp int64) error {

	stmt, err := r.Conn.Db.Prepare(`
	INSERT INTO 
		Jail(
			userId, 
			reason,
			task,
			jailedAt
		)
	VALUES(?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, reason, task, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (r JailRepository) RemoveUserFromJail(userId string) error {

	query := "DELETE FROM Jail WHERE userId = ?"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting jail entry for user: %w", err)
	}

	return nil
}

func (r JailRepository) GetJailedUser(userId string) (*dataModels.JailedUser, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Jail WHERE userId = ?"
	row := r.Conn.Db.QueryRow(query, userId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var jailedUser dataModels.JailedUser
	err := row.Scan(&jailedUser.UserId,
		&jailedUser.Reason,
		&jailedUser.TaskToComplete,
		&jailedUser.JailedAt,
	)

	if err != nil {
		return nil, err
	}

	return &jailedUser, nil

}
