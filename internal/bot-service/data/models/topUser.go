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

type TopUserADS struct {
	DiscordTag string
	UserId     string
	Streak     int // in days
}

type TopUserRCT struct {
	DiscordTag        string
	UserId            string
	ReactionsReceived int
}
