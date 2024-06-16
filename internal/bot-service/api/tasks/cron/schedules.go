package cron

import "time"

// Returns a delay and a ticker to use for the initial delay
// and then subsequent executions of the activity streak update cron.
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

// GetDelayAndTickerForMonthlyLeaderboardCron returns a delay and a ticker to use for
// the initial delay and then subsequent executions of the monthly leaderboard update cron.
func GetDelayAndTickerForMonthlyLeaderboardCron(lastDay bool, hourInDay int, minuteInDay int, secondInDay int) (time.Duration, *time.Ticker) {

	now := time.Now()

	// Get last day of current month
	year, month, _ := now.Date()
	nextMonth := time.Date(year, month+1, 1, hourInDay, minuteInDay, secondInDay, 0, now.Location())
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)

	// TODO: Remove this check once the feature is properly tested
	if !lastDay {
		lastDayOfMonth = time.Date(year, month, now.Day(), hourInDay, minuteInDay, secondInDay, 0, now.Location())
	}

	// Calculate the duration until the next target hour
	nextRun := time.Date(year, month, lastDayOfMonth.Day(), lastDayOfMonth.Hour(), lastDayOfMonth.Minute(), lastDayOfMonth.Second(), 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.AddDate(0, 1, 0) // Move to the next exec time if the target time has passed today
	}

	oneMonthDuration := time.Hour * 24 * 30 // roughly 30 days in one month
	return nextRun.Sub(now), time.NewTicker(oneMonthDuration)

}
