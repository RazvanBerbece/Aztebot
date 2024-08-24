package events

type CoinAwardEvent struct {
	GuildId  string
	UserId   string
	Funds    float64
	Activity string
}
