package cron

import "time"

// Returns a delay and a ticker to use for the initial delay and then subsequent executions of the activity streak update cron.
func GetDelayAndTickerForActivityStreakCron(hour int, minute int, second int) (time.Duration, *time.Ticker) {

	// Run activity streak logic at given timestamp
	targetHour := hour
	targetMinute := minute
	targetSecond := second

	// Calculate the duration until the next target hour
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, targetSecond, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour) // Move to the next day if the target hour has passed today
	}

	return nextRun.Sub(now), time.NewTicker(time.Hour * 24)

}
