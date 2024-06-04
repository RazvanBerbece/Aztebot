package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type ArcadeLadderRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewArcadeLadderRepository() *ArcadeLadderRepository {
	repo := new(ArcadeLadderRepository)
	repo.Conn.Connect()
	return repo
}

func (r ArcadeLadderRepository) EntryExists(userId string) int {
	query := "SELECT COUNT(*) FROM ArcadeLadder WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for arcade ladder entry: %v\n", err)
		return -1
	}
	return count
}

func (r ArcadeLadderRepository) AddNewLadderEntry(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
	INSERT INTO 
	ArcadeLadder(
			userId, 
			wins
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

func (r ArcadeLadderRepository) DeleteEntry(userId string) error {

	query := "DELETE FROM ArcadeLadder WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting arcade ladder entry: %w", err)
	}

	return nil
}

func (r ArcadeLadderRepository) AddWin(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE ArcadeLadder SET 
			wins = wins + ?
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

func (r ArcadeLadderRepository) RemoveWin(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE ArcadeLadder SET 
			wins = wins - ?
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

func (r ArcadeLadderRepository) GetArcadeLadder() ([]dax.ArcadeLadderEntry, error) {

	var entries []dax.ArcadeLadderEntry

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM ArcadeLadder ORDER BY wins DESC")
	if err != nil {
		return nil, fmt.Errorf("GetArcadeLadder: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry dax.ArcadeLadderEntry
		if err := rows.Scan(&entry.UserId, &entry.Wins); err != nil {
			return nil, fmt.Errorf("GetArcadeLadder: %v", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetArcadeLadder: %v", err)
	}

	return entries, nil

}
