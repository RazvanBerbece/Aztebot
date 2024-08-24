package dax

type CurrencySystemState struct {
	GuildId                string
	CurrencyName           string
	TotalCurrencyAvailable float64
	TotalCurrencyInFlow    float64
	DateOfLastReplenish    int64
}
