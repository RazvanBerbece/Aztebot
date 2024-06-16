package dataModels

type TopUserMS struct {
	DiscordTag   string
	UserId       string
	MessagesSent int
}

type TopUserVC struct {
	DiscordTag     string
	UserId         string
	TimeSpentInVCs int // in seconds
}
