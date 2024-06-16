package cron

import (
	"database/sql"
	"fmt"
	"time"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

func ProcessUpdateActivityStreaks(h int, m int, s int) {

	initialActivityStreakDelay, activityStreakTicker := GetDelayAndTickerForActivityStreakCron(h, m, s) // H, m, s

	go func() {

		fmt.Println("[SCHEDULED CRON] Scheduled Task UpdateActivityStreaks() in <", initialActivityStreakDelay.Hours(), "> hours")
		time.Sleep(initialActivityStreakDelay)

		// Inject new connections
		usersRepository := repositories.NewUsersRepository()
		userStatsRepository := repositories.NewUsersStatsRepository()

		// The first run should happen at start-up, not after 24 hours
		UpdateActivityStreaks(globalRepositories.UsersRepository, globalRepositories.UserStatsRepository)

		for range activityStreakTicker.C {
			// Process
			UpdateActivityStreaks(usersRepository, userStatsRepository)
		}
	}()
}

func UpdateActivityStreaks(usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) {

	fmt.Println("[CRON] Starting Task UpdateActivityStreaks() at", time.Now())

	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("[STARTUP] Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
	}

	// For all users in the database
	fmt.Println("[CRON] Checkpoint Task UpdateActivityStreaks() at", time.Now(), "-> Updating", len(uids), "streaks")
	for _, uid := range uids {
		stats, err := userStatsRepository.GetStatsForUser(uid)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("[CRON] Failed Task UpdateActivityStreaks() at", time.Now(), "for UID", "with error", err)
		}

		// lastActiveSince smaller than 24 (which means did an action in the last 24 hours)
		timestampTime := time.Unix(stats.LastActiveTimestamp, 0)
		lastActiveSince := time.Since(timestampTime)

		// Activity scores greater than this are favourable
		var activityThreshold int
		if globalConfiguration.FavourableActivitiesThresholdErr != nil {
			activityThreshold = 10
		} else {
			activityThreshold = globalConfiguration.FavourableActivitiesThreshold
		}

		// If user has favourable activity score and favourable timestamp, increase day streak
		if lastActiveSince.Hours() < 24 && stats.NumberActivitiesToday > activityThreshold {
			err := userStatsRepository.IncrementActiveDayStreakForUser(uid)
			if err != nil {
				fmt.Println("[CRON] Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
			}
		} else {
			err := userStatsRepository.ResetActiveDayStreakForUser(uid)
			if err != nil {
				fmt.Println("[CRON] Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
			}
		}

		// Reset the activity count for the next day
		err = userStatsRepository.ResetActivitiesTodayForUser(uid)
		if err != nil {
			fmt.Println("[CRON] Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
		}
	}

	fmt.Println("[CRON] Finished Task UpdateActivityStreaks() at", time.Now())

}
