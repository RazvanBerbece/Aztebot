package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type MonthlyLeaderboardRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewMonthlyLeaderboardRepository() *MonthlyLeaderboardRepository {
	repo := new(MonthlyLeaderboardRepository)
	repo.Conn.Connect()
	return repo
}

func (r MonthlyLeaderboardRepository) EntryExists(userId string) int {
	query := "SELECT COUNT(*) FROM MonthlyLeaderboard WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking monthly entry: %v\n", err)
		return -1
	}
	return count
}

func (r MonthlyLeaderboardRepository) AddLeaderboardEntry(userId string, category int8) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
	INSERT INTO 
		MonthlyLeaderboard(
			userId, 
			xpEarnedInCurrentMonth, 
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

func (r MonthlyLeaderboardRepository) DeleteEntry(userId string) error {

	query := "DELETE FROM MonthlyLeaderboard WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting monthly leaderboard entry: %w", err)
	}

	return nil
}

func (r MonthlyLeaderboardRepository) ClearLeaderboard() error {

	query := "TRUNCATE TABLE MonthlyLeaderboard"

	_, err := r.Conn.SqlDb.Exec(query)
	if err != nil {
		return fmt.Errorf("error clearing monthly leaderboard: %w", err)
	}

	return nil
}

func (r MonthlyLeaderboardRepository) AddLeaderboardExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE MonthlyLeaderboard SET 
			xpEarnedInCurrentMonth = xpEarnedInCurrentMonth + ?
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

func (r MonthlyLeaderboardRepository) RemoveUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE MonthlyLeaderboard SET 
			xpEarnedInCurrentMonth = xpEarnedInCurrentMonth - ?
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

func (r MonthlyLeaderboardRepository) GetLeaderboardEntriesByCategory(category int8) ([]dax.MonthlyLeaderboardEntry, error) {

	var entries []dax.MonthlyLeaderboardEntry

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM MonthlyLeaderboard WHERE category = ? AND xpEarnedInCurrentMonth > 0 ORDER BY xpEarnedInCurrentMonth DESC", category)
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

func (r MonthlyLeaderboardRepository) UpdateCategoryForUser(userId string, category int8) error {
	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE MonthlyLeaderboard SET 
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
