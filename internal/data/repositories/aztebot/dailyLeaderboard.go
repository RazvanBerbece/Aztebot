package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type DailyLeaderboardRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewDailyLeaderboardRepository() *DailyLeaderboardRepository {
	repo := new(DailyLeaderboardRepository)
	repo.Conn.Connect()
	return repo
}

func (r DailyLeaderboardRepository) EntryExists(userId string) int {
	query := "SELECT COUNT(*) FROM DailyLeaderboard WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for daily entry: %v\n", err)
		return -1
	}
	return count
}

func (r DailyLeaderboardRepository) AddLeaderboardEntry(userId string, category int8) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
	INSERT INTO 
		DailyLeaderboard(
			userId, 
			xpEarnedInCurrentDay, 
			category
		)
	VALUES(?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, 0, category)
	if err != nil {
		return err
	}

	return nil
}

func (r DailyLeaderboardRepository) DeleteEntry(userId string) error {

	query := "DELETE FROM DailyLeaderboard WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting daily leaderboard entry: %w", err)
	}

	return nil
}

func (r DailyLeaderboardRepository) ClearLeaderboard() error {

	query := "TRUNCATE TABLE DailyLeaderboard"

	_, err := r.Conn.SqlDb.Exec(query)
	if err != nil {
		return fmt.Errorf("error clearing daily leaderboard: %w", err)
	}

	return nil
}

func (r DailyLeaderboardRepository) AddLeaderboardExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE DailyLeaderboard SET 
			xpEarnedInCurrentDay = xpEarnedInCurrentDay + ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(experiencePoints, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r DailyLeaderboardRepository) RemoveUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE DailyLeaderboard SET 
			xpEarnedInCurrentDay = xpEarnedInCurrentDay - ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(experiencePoints, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r DailyLeaderboardRepository) GetLeaderboardEntriesByCategory(category int8) ([]dax.MonthlyLeaderboardEntry, error) {

	var entries []dax.MonthlyLeaderboardEntry

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM DailyLeaderboard WHERE category = ? AND xpEarnedInCurrentDay > 0 ORDER BY xpEarnedInCurrentDay DESC", category)
	if err != nil {
		return nil, fmt.Errorf("GetAllLeaderboardEntries: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry dax.MonthlyLeaderboardEntry
		if err := rows.Scan(&entry.UserId, &entry.XpEarnedInCurrentMonth, &entry.Category); err != nil {
			return nil, fmt.Errorf("GetAllLeaderboardEntries: %v", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllLeaderboardEntries: %v", err)
	}

	return entries, nil

}

func (r DailyLeaderboardRepository) UpdateCategoryForUser(userId string, category int8) error {
	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE DailyLeaderboard SET 
			category = ?
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing to update a user's leaderboard category %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category, userId)
	if err != nil {
		fmt.Printf("Error ocurred while updating leaderboard category for user %s: %v\n", userId, err)
		return err
	}

	return nil
}
