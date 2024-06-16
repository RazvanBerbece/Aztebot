package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func ProcessRemoveExpiredWarns(months int) {

	initialWarnRemovalDelay, warnRemovalTicker := GetDelayAndTickerForWarnRemovalCron(months) // every n=2 months

	go func() {

		RemoveExpiredWarns(globalsRepo.WarnsRepository)

		fmt.Println("[SCHEDULED CRON]\t\tScheduled Task RemoveExpiredWarns() in <", initialWarnRemovalDelay.Hours()/24, "> days")
		time.Sleep(initialWarnRemovalDelay)

		// Inject new connections
		warnsRepository := repositories.NewWarnsRepository()

		for range warnRemovalTicker.C {
			// Process
			RemoveExpiredWarns(warnsRepository)
		}

		// Cleanup DB connections after cron run
		utils.CleanupRepositories(nil, nil, nil, warnsRepository, nil)
	}()
}

func RemoveExpiredWarns(warnsRepository *repositories.WarnsRepository) {

	fmt.Println("[CRON]\t\t\tStarting Task RemoveExpiredWarns() at", time.Now())

	allWarns, err := warnsRepository.GetAllWarns()
	if err != nil {
		fmt.Println("[CRON]\t\t\tFailed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
	}

	// For all existing warns
	for _, warn := range allWarns {
		warnCreationTime := time.Unix(warn.CreationTimestamp, 0)
		twoMonthsAgo := time.Hour * 24 * 61
		// If the warn is older than 2 months
		if time.Since(warnCreationTime) > twoMonthsAgo {
			// Remove it
			err := warnsRepository.DeleteWarningForUser(warn.Id, warn.UserId)
			if err != nil {
				fmt.Println("[CRON]\t\t\tFailed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
			}
		}
	}

	fmt.Println("[CRON]\t\t\tFinished Task RemoveExpiredWarns() at", time.Now())

}
