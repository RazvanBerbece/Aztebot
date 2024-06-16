package repositories

import (
	"database/sql"
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type UserRepRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewUserRepRepository() *UserRepRepository {
	repo := new(UserRepRepository)
	repo.Conn.Connect()
	return repo
}

func (r UserRepRepository) EntryExists(userId string) int {
	query := "SELECT COUNT(*) FROM UserRep WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for user rep entry: %v\n", err)
		return -1
	}
	return count
}

func (r UserRepRepository) AddNewEntry(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
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

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user rep entry: %w", err)
	}

	return nil
}

func (r UserRepRepository) AddRep(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
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

	stmt, err := r.Conn.SqlDb.Prepare(`
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

func (r UserRepRepository) ResetRep(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE UserRep SET 
			rep = 0
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepRepository) GetRepForUser(userId string) (*dax.UserRep, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM UserRep WHERE userId = ?"
	row := r.Conn.SqlDb.QueryRow(query, userId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var userRep dax.UserRep
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

func (r UserRepRepository) GetRepTop() ([]dax.UserRep, error) {

	var entries []dax.UserRep

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM UserRep ORDER BY rep DESC")
	if err != nil {
		return nil, fmt.Errorf("GetRepTop: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry dax.UserRep
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
