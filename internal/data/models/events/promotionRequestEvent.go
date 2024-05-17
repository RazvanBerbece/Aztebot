package events

type PromotionRequestEvent struct {
	GuildId       string
	UserId        string
	UserTag       string
	CurrentXp     float64
	MessagesSent  int
	TimeSpentInVc int
	CurrentLevel  int
}
