package repositories

import (
	"fmt"
	"time"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type UsersStatsRepository struct {
	Conn databaseconn.Database
}

func NewUsersStatsRepository() *UsersStatsRepository {
	repo := new(UsersStatsRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r UsersStatsRepository) UserStatsExist(userId string) int {
	query := "SELECT COUNT(*) FROM Users WHERE userId = ?"
	var count int
	err := r.Conn.Db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r UsersStatsRepository) SaveInitialUserStats(userId string) error {

	userStats := &dataModels.UserStats{
		UserId:                    userId,
		NumberMessagesSent:        0,
		NumberSlashCommandsUsed:   0,
		NumberReactionsReceived:   0,
		NumberActiveDayStreak:     0,
		LastActiveTimestamp:       0,
		NumberActivitiesToday:     1,
		TimeSpentInVoiceChannels:  0,
		TimeSpentInEvents:         0,
		TimeSpentListeningToMusic: 0,
	}

	stmt, err := r.Conn.Db.Prepare(`
		INSERT INTO 
			UserStats(
				userId, 
				messagesSent, 
				slashCommandsUsed, 
				reactionsReceived, 
				activeDayStreak,
				lastActivityTimestamp,
				numberActivitiesToday,
				timeSpentInVoiceChannels,
				timeSpentInEvents,
				timeSpentListeningMusic
			)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userStats.UserId, userStats.NumberMessagesSent, userStats.NumberSlashCommandsUsed, userStats.NumberReactionsReceived, userStats.NumberActiveDayStreak, userStats.LastActiveTimestamp, userStats.NumberActivitiesToday, userStats.TimeSpentInVoiceChannels, userStats.TimeSpentInEvents, userStats.TimeSpentListeningToMusic)
	if err != nil {
		return err
	}

	return nil

}

func (r UsersStatsRepository) GetStatsForUser(userId string) (*dataModels.UserStats, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM UserStats WHERE userId = ?"
	row := r.Conn.Db.QueryRow(query, userId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var userStats dataModels.UserStats
	err := row.Scan(&userStats.Id,
		&userStats.UserId,
		&userStats.NumberMessagesSent,
		&userStats.NumberSlashCommandsUsed,
		&userStats.NumberReactionsReceived,
		&userStats.NumberActiveDayStreak,
		&userStats.LastActiveTimestamp,
		&userStats.NumberActivitiesToday,
		&userStats.TimeSpentInVoiceChannels,
		&userStats.TimeSpentInEvents,
		&userStats.TimeSpentListeningToMusic,
	)

	if err != nil {
		return nil, err
	}

	return &userStats, nil

}

func (r UsersStatsRepository) DeleteUserStats(userId string) error {

	query := "DELETE FROM UserStats WHERE userId = ?"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user stats: %w", err)
	}

	return nil
}

func (r UsersStatsRepository) IncrementMessagesSentForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			messagesSent = messagesSent + 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing messages sent stat increment for user %s: %v", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while incrementing messages sent stat for user %s: %v", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) DecrementMessagesSentForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			messagesSent = messagesSent - 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing messages sent stat decrement for user %s: %v", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while decrementing messages sent stat for user %s: %v", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) IncrementSlashCommandsUsedForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			slashCommandsUsed = slashCommandsUsed + 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing slash commands used stat increment for user %s: %v", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while incrementing slash commands used stat for user %s: %v", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) IncrementReactionsReceivedForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			reactionsReceived = reactionsReceived + 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing reactions received stat increment for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while incrementing reactions received stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) DecrementReactionsReceivedForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			reactionsReceived = reactionsReceived - 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing reactions received stat decrement for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while decrementing reactions received stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) IncrementActiveDayStreakForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			activeDayStreak = activeDayStreak + 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing active day streak stat increment for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while incrementing active day streak stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) ResetActiveDayStreakForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			activeDayStreak = 0
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing active day streak stat reset for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	retries := 5
	for i := 0; i < retries; i++ {
		_, err = stmt.Exec(userId)
		if err != nil {
			fmt.Printf("Error ocurred while resetting day streak stat for user %s: %v\nRetrying...", userId, err)
			time.Sleep(time.Millisecond * 20)
		} else {
			break
		}
	}

	return nil
}

func (r UsersStatsRepository) UpdateLastActiveTimestamp(userId string, timestamp int64) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			lastActivityTimestamp = ?
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing last active timestamp for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(timestamp, userId)
	if err != nil {
		fmt.Printf("Error ocurred while updating the last active timestamp for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) IncrementActivitiesTodayForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			numberActivitiesToday = numberActivitiesToday + 1
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing activities number increment for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Printf("Error ocurred while incrementing activities number stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) ResetActivitiesTodayForUser(userId string) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			numberActivitiesToday = 0
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing reset activities number for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	retries := 5
	for i := 0; i < retries; i++ {
		_, err = stmt.Exec(userId)
		if err != nil {
			fmt.Printf("Error ocurred while resetting activities number stat for user %s: %v\nRetrying...", userId, err)
			time.Sleep(time.Millisecond * 20)
		} else {
			break
		}
	}

	return nil
}

func (r UsersStatsRepository) AddToTimeSpentInVoiceChannels(userId string, sTimeLength int) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			timeSpentInVoiceChannels = timeSpentInVoiceChannels + ?
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing VC spent time increase for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sTimeLength, userId)
	if err != nil {
		fmt.Printf("Error ocurred while increasing VC spent time stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) AddToTimeSpentInEvents(userId string, sTimeLength int) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			timeSpentInEvents = timeSpentInEvents + ?
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing event spent time increase for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, sTimeLength)
	if err != nil {
		fmt.Printf("Error ocurred while increasing event spent time stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) AddToTimeSpentListeningMusic(userId string, sTimeLength int) error {
	stmt, err := r.Conn.Db.Prepare(`
		UPDATE UserStats SET 
			timeSpentListeningMusic = timeSpentListeningMusic + ?
		WHERE userId = ?`)
	if err != nil {
		fmt.Printf("Error ocurred while preparing music listening spent time increase for user %s: %v\n", userId, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sTimeLength, userId)
	if err != nil {
		fmt.Printf("Error ocurred while increasing music listening spent time stat for user %s: %v\n", userId, err)
		return err
	}

	return nil
}

func (r UsersStatsRepository) GetTopUsersByMessageSent(count int) ([]dataModels.TopUserMS, error) {

	// This could use something similar to a strategy pattern
	// and only pass the column we want to filter on as a parameter to a more generic function

	query := `SELECT Users.discordTag, UserStats.userId, UserStats.messagesSent
		FROM UserStats
		JOIN Users ON UserStats.userId = Users.userId
		WHERE UserStats.messagesSent > 10
		ORDER BY UserStats.messagesSent DESC
		LIMIT ?`

	rows, err := r.Conn.Db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []dataModels.TopUserMS

	for rows.Next() {
		var user dataModels.TopUserMS
		err := rows.Scan(&user.DiscordTag, &user.UserId, &user.MessagesSent)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topUsers, nil

}

func (r UsersStatsRepository) GetTopUsersByTimeSpentInVC(count int) ([]dataModels.TopUserVC, error) {

	// This could use something similar to a strategy pattern
	// and only pass the column we want to filter on as a parameter to a more generic function

	query := `SELECT Users.discordTag, UserStats.userId, UserStats.timeSpentInVoiceChannels
		FROM UserStats
		JOIN Users ON UserStats.userId = Users.userId
		WHERE UserStats.timeSpentInVoiceChannels > 5
		ORDER BY UserStats.timeSpentInVoiceChannels DESC
		LIMIT ?`

	rows, err := r.Conn.Db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []dataModels.TopUserVC

	for rows.Next() {
		var user dataModels.TopUserVC
		err := rows.Scan(&user.DiscordTag, &user.UserId, &user.TimeSpentInVCs)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topUsers, nil

}

func (r UsersStatsRepository) GetTopUsersByTimeSpentListeningMusic(count int) ([]dataModels.TopUserMusic, error) {

	// This could use something similar to a strategy pattern
	// and only pass the column we want to filter on as a parameter to a more generic function

	query := `SELECT Users.discordTag, UserStats.userId, UserStats.timeSpentListeningMusic
		FROM UserStats
		JOIN Users ON UserStats.userId = Users.userId
		WHERE UserStats.timeSpentListeningMusic > 10
		ORDER BY UserStats.timeSpentListeningMusic DESC
		LIMIT ?`

	rows, err := r.Conn.Db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []dataModels.TopUserMusic

	for rows.Next() {
		var user dataModels.TopUserMusic
		err := rows.Scan(&user.DiscordTag, &user.UserId, &user.TimeSpentListeningMusic)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topUsers, nil

}

func (r UsersStatsRepository) GetTopUsersByActiveDayStreak(count int) ([]dataModels.TopUserADS, error) {

	// This could use something similar to a strategy pattern
	// and only pass the column we want to filter on as a parameter to a more generic function

	query := `SELECT Users.discordTag, UserStats.userId, UserStats.activeDayStreak
		FROM UserStats
		JOIN Users ON UserStats.userId = Users.userId
		WHERE UserStats.activeDayStreak > 1
		ORDER BY UserStats.activeDayStreak DESC
		LIMIT ?`

	rows, err := r.Conn.Db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []dataModels.TopUserADS

	for rows.Next() {
		var user dataModels.TopUserADS
		err := rows.Scan(&user.DiscordTag, &user.UserId, &user.Streak)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topUsers, nil

}

func (r UsersStatsRepository) GetTopUsersByReceivedReactions(count int) ([]dataModels.TopUserRCT, error) {

	// This could use something similar to a strategy pattern
	// and only pass the column we want to filter on as a parameter to a more generic function

	query := `SELECT Users.discordTag, UserStats.userId, UserStats.reactionsReceived
		FROM UserStats
		JOIN Users ON UserStats.userId = Users.userId
		WHERE UserStats.reactionsReceived > 1
		ORDER BY UserStats.reactionsReceived DESC
		LIMIT ?`

	rows, err := r.Conn.Db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []dataModels.TopUserRCT

	for rows.Next() {
		var user dataModels.TopUserRCT
		err := rows.Scan(&user.DiscordTag, &user.UserId, &user.ReactionsReceived)
		if err != nil {
			return nil, err
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topUsers, nil

}
