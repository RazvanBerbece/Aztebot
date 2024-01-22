package cronFeature

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/cron"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func ProcessRemoveExpiredWarns(months int) {
	initialWarnRemovalDelay, warnRemovalTicker := cron.GetDelayAndTickerForWarnRemovalCron(months) // every n=2 months
	go func() {

		fmt.Println("Scheduled Task RemoveExpiredWarns() in <", initialWarnRemovalDelay.Hours()/24, "> days")
		time.Sleep(initialWarnRemovalDelay)

		RemoveExpiredWarns(globalsRepo.WarnsRepository)

		for range warnRemovalTicker.C {
			// Inject new connections
			warnsRepository := repositories.NewWarnsRepository()

			// Process
			RemoveExpiredWarns(warnsRepository)

			// Cleanup DB connections after cron run
			utils.CleanupRepositories(nil, nil, nil, warnsRepository)
		}
	}()
}

func RemoveExpiredWarns(warnsRepository *repositories.WarnsRepository) {

	fmt.Println("Starting Task RemoveExpiredWarns() at", time.Now())

	allWarns, err := warnsRepository.GetAllWarns()
	if err != nil {
		fmt.Println("Failed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
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
				fmt.Println("Failed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
			}
		}
	}

	fmt.Println("Finished Task RemoveExpiredWarns() at", time.Now())

}
