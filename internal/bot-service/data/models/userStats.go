package dataModels

type UserStats struct {
	Id                       int
	UserId                   string
	NumberMessagesSent       int
	NumberSlashCommandsUsed  int
	NumberReactionsReceived  int
	NumberActiveDayStreak    int
	LastActiveTimestamp      int64
	NumberActivitiesToday    int
	TimeSpentInVoiceChannels int
	TimeSpentInEvents        int
}
