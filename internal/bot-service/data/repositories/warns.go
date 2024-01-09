package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type WarnsRepository struct {
	Conn databaseconn.Database
}

func NewWarnsRepository() *WarnsRepository {
	repo := new(WarnsRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r WarnsRepository) SaveWarn(userId string, reason string, timestamp int64) error {

	userStats := &dataModels.Warn{
		UserId:    userId,
		Reason:    reason,
		Timestamp: timestamp,
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
