package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type JailRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewJailRepository() *JailRepository {
	repo := new(JailRepository)
	repo.Conn.Connect()
	return repo
}

func (r JailRepository) UserIsJailed(userId string) int {
	query := "SELECT COUNT(*) FROM Jail WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for jail entry: %v\n", err)
		return -1
	}
	return count
}

func (r JailRepository) AddUserToJail(userId string, reason string, task string, timestamp int64, roleIdsBeforeJail string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
	INSERT INTO 
		Jail(
			userId, 
			reason,
			task,
			jailedAt,
			roleIdsBeforeJail
		)
	VALUES(?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, reason, task, timestamp, roleIdsBeforeJail)
	if err != nil {
		return err
	}

	return nil
}

func (r JailRepository) RemoveUserFromJail(userId string) error {

	query := "DELETE FROM Jail WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting jail entry for user: %w", err)
	}

	return nil
}

func (r JailRepository) GetJailedUser(userId string) (*dax.JailedUser, error) {

	query := "SELECT * FROM Jail WHERE userId = ?"
	row := r.Conn.SqlDb.QueryRow(query, userId)

	var jailedUser dax.JailedUser
	err := row.Scan(&jailedUser.UserId,
		&jailedUser.Reason,
		&jailedUser.TaskToComplete,
		&jailedUser.JailedAt,
		&jailedUser.RoleIdsBeforeJail,
	)

	if err != nil {
		return nil, err
	}

	return &jailedUser, nil

}

func (r JailRepository) GetJail() ([]dax.JailedUser, error) {

	var jailed []dax.JailedUser

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM Jail ORDER BY jailedAt ASC")
	if err != nil {
		return nil, fmt.Errorf("GetJail: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var jailedUser dax.JailedUser
		if err := rows.Scan(&jailedUser.UserId, &jailedUser.Reason, &jailedUser.TaskToComplete, &jailedUser.JailedAt, &jailedUser.RoleIdsBeforeJail); err != nil {
			return nil, fmt.Errorf("GetJail: %v", err)
		}
		jailed = append(jailed, jailedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetJail: %v", err)
	}

	return jailed, nil

}
