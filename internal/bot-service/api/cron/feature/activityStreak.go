package cronFeature

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/cron"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func ProcessUpdateActivityStreaks(h int, m int, s int) {

	initialActivityStreakDelay, activityStreakTicker := cron.GetDelayAndTickerForActivityStreakCron(h, m, s) // H, m, s

	go func() {

		fmt.Println("Scheduled Task UpdateActivityStreaks() in <", initialActivityStreakDelay.Hours(), "> hours")
		time.Sleep(initialActivityStreakDelay)

		// The first run should happen at start-up, not after 24 hours
		UpdateActivityStreaks(globalsRepo.UsersRepository, globalsRepo.UserStatsRepository)

		for range activityStreakTicker.C {
			// Inject new connections
			usersRepository := repositories.NewUsersRepository()
			userStatsRepository := repositories.NewUsersStatsRepository()

			// Process
			UpdateActivityStreaks(usersRepository, userStatsRepository)

			// Cleanup DB connections after cron run
			utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil)
		}
	}()
}

func UpdateActivityStreaks(usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) {

	fmt.Println("Starting Task UpdateActivityStreaks() at", time.Now())

	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
	}

	// For all users in the database
	fmt.Println("Checkpoint Task UpdateActivityStreaks() at", time.Now(), "-> Updating", len(uids), "streaks")
	for _, uid := range uids {
		stats, err := userStatsRepository.GetStatsForUser(uid)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("Failed Task UpdateActivityStreaks() at", time.Now(), "for UID", "with error", err)
		}

		// lastActiveSince smaller than 24 (which means did an action in the last 24 hours)
		timestampTime := time.Unix(stats.LastActiveTimestamp, 0)
		lastActiveSince := time.Since(timestampTime)

		// Activity scores greater than this are favourable
		var activityThreshold int
		if globals.FavourableActivitiesThresholdErr != nil {
			activityThreshold = 10
		} else {
			activityThreshold = globals.FavourableActivitiesThreshold
		}

		// If user has favourable activity score and favourable timestamp, increase day streak
		if lastActiveSince.Hours() < 24 && stats.NumberActivitiesToday > activityThreshold {
			err := userStatsRepository.IncrementActiveDayStreakForUser(uid)
			if err != nil {
				fmt.Println("Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
			}
		} else {
			err := userStatsRepository.ResetActiveDayStreakForUser(uid)
			if err != nil {
				fmt.Println("Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
			}
		}

		// Reset the activity count for the next day
		err = userStatsRepository.ResetActivitiesTodayForUser(uid)
		if err != nil {
			fmt.Println("Failed Task UpdateActivityStreaks() at", time.Now(), "with error", err)
		}
	}

	fmt.Println("Finished Task UpdateActivityStreaks() at", time.Now())

}
