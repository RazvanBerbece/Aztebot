package dax

type MonthlyLeaderboardEntry = struct {
	UserId                 string
	XpEarnedInCurrentMonth float64
	Category               int8
}
