package repositories

import (
	"database/sql"
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
)

type UserRepRepository struct {
	Conn databaseconn.Database
}

func NewUserRepRepository() *UserRepRepository {
	repo := new(UserRepRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r UserRepRepository) EntryExists(userId string) int {
	query := "SELECT COUNT(*) FROM UserRep WHERE userId = ?"
	var count int
	err := r.Conn.Db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for user rep entry: %v\n", err)
		return -1
	}
	return count
}

func (r UserRepRepository) AddNewEntry(userId string) error {

	stmt, err := r.Conn.Db.Prepare(`
	INSERT INTO 
	UserRep(
			userId, 
			rep
		)
	VALUES(?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, 0)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepRepository) DeleteEntry(userId string) error {

	query := "DELETE FROM UserRep WHERE userId = ?"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user rep entry: %w", err)
	}

	return nil
}

func (r UserRepRepository) AddRep(userId string) error {

	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserRep SET 
			rep = rep + ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepRepository) RemoveRep(userId string) error {

	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserRep SET 
			rep = rep - ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepRepository) GetRepForUser(userId string) (*dataModels.UserRep, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM UserRep WHERE userId = ?"
	row := r.Conn.Db.QueryRow(query, userId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var userRep dataModels.UserRep
	err := row.Scan(&userRep.UserId, &userRep.Rep)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &userRep, nil

}

func (r UserRepRepository) GetRepTop() ([]dataModels.UserRep, error) {

	var entries []dataModels.UserRep

	rows, err := r.Conn.Db.Query("SELECT * FROM UserRep ORDER BY rep DESC")
	if err != nil {
		return nil, fmt.Errorf("GetRepTop: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry dataModels.UserRep
		if err := rows.Scan(&entry.UserId, &entry.Rep); err != nil {
			return nil, fmt.Errorf("GetRepTop: %v", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRepTop: %v", err)
	}

	return entries, nil

}
