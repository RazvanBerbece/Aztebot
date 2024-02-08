package cron

import "time"

// Returns a delay and a ticker to use for the initial delay and then subsequent executions of the activity streak update cron.
func GetDelayAndTickerForActivityStreakCron(hour int, minute int, second int) (time.Duration, *time.Ticker) {

	// Run ativity streak logic at given timestamp
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

// Returns a delay and a ticker to use for the initial delay and then subsequent executions of the warn removal cron.
func GetDelayAndTickerForWarnRemovalCron(months int) (time.Duration, *time.Ticker) {

	now := time.Now()

	// Run cron logic at given timestamp
	targetMonth := now.AddDate(0, months, 0).Month()

	// Calculate the duration until the next target hour
	nextRun := time.Date(now.Year(), targetMonth, now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.AddDate(0, 2, 0) // Move to the next exec time if the target time has passed today
	}

	twoMonthsDuration := time.Hour * 24 * 61 // roughly 61 days in 2 months, 30 and 31 etc.
	return nextRun.Sub(now), time.NewTicker(twoMonthsDuration)

}

// Returns a delay and a ticker to use for the initial delay and then subsequent executions of the archived timeout cleanup cron.
func GetDelayAndTickerForArchiveCleanupCron(months int) (time.Duration, *time.Ticker) {

	now := time.Now()

	// Run cron logic at given timestamp
	targetMonth := now.AddDate(0, months, 0).Month()

	// Calculate the duration until the next target hour
	nextRun := time.Date(now.Year(), targetMonth, now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.AddDate(0, 2, 0) // Move to the next exec time if the target time has passed today
	}

	oneMonthDuration := time.Hour * 24 * 30
	return nextRun.Sub(now), time.NewTicker(oneMonthDuration)

}
