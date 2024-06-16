package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func ProcessRemoveArchivedTimeouts(months int) {

	initialWarnRemovalDelay, warnRemovalTicker := GetDelayAndTickerForWarnRemovalCron(months)

	go func() {

		CleanupArchivedTimeouts(globalsRepo.TimeoutsRepository)

		fmt.Println("[SCHEDULED CRON] Scheduled Task ProcessRemoveArchivedTimeouts() in <", initialWarnRemovalDelay.Hours()/24, "> days")
		time.Sleep(initialWarnRemovalDelay)

		// Inject new connections
		timeoutsRepository := repositories.NewTimeoutsRepository()

		for range warnRemovalTicker.C {
			// Process
			CleanupArchivedTimeouts(timeoutsRepository)
		}

		// Cleanup DB connections after cron run
		utils.CleanupRepositories(nil, nil, nil, nil, timeoutsRepository)
	}()
}

func CleanupArchivedTimeouts(timeoutsRepository *repositories.TimeoutsRepository) {

	fmt.Println("[CRON] Starting Task CleanupArchivedTimeouts() at", time.Now())

	archivedTimeouts, err := timeoutsRepository.GetAllArchivedTimeouts()
	if err != nil {
		fmt.Println("[CRON]Failed Task CleanupArchivedTimeouts() at", time.Now(), "with error", err)
	}

	// For all archived timeouts
	for _, timeout := range archivedTimeouts {
		timeoutExpiryTime := time.Unix(timeout.ExpiryTimestamp, 0)
		aMonthAgo := time.Hour * 24 * 30
		// If the archived timeout is older than 1 month
		if time.Since(timeoutExpiryTime) > aMonthAgo {
			// Remove it
			err := timeoutsRepository.ClearArchivedTimeout(timeout.Id)
			if err != nil {
				fmt.Println("[CRON]Failed Task CleanupArchivedTimeouts() at", time.Now(), "with error", err)
			}
		}
	}

	fmt.Println("[CRON] Finished Task CleanupArchivedTimeouts() at", time.Now())

}
