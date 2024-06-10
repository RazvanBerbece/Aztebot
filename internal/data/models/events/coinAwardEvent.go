package events

type CoinAwardEvent struct {
	UserId   string
	Funds    float64
	Activity string
}
