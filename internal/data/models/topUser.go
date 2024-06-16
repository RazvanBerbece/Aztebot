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

type TopUserMusic struct {
	DiscordTag              string
	UserId                  string
	TimeSpentListeningMusic int // in seconds
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

type TopUserXP struct {
	DiscordTag string
	UserId     string
	XpGained   float64
}
