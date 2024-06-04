package cron

import (
	"fmt"
	"time"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
)

func ProcessRemoveExpiredWarns() {

	// TODO: Move this into env variable and use cron syntax
	var numSec int = 5.256e+6 // run every 2 months

	fmt.Println("[CRON] Starting Cron Ticker ProcessRemoveExpiredWarns() at", time.Now(), "running every 2 months")

	// Inject new connections
	warnsRepository := repositories.NewWarnsRepository()

	go RemoveExpiredWarns(warnsRepository)

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go RemoveExpiredWarns(warnsRepository)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func RemoveExpiredWarns(warnsRepository *repositories.WarnsRepository) {

	fmt.Println("[CRON] Starting Task RemoveExpiredWarns() at", time.Now())

	allWarns, err := warnsRepository.GetAllWarns()
	if err != nil {
		fmt.Println("[CRON] Failed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
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
				fmt.Println("[CRON] Failed Task RemoveExpiredWarns() at", time.Now(), "with error", err)
			}
		}
	}

	fmt.Println("[CRON] Finished Task RemoveExpiredWarns() at", time.Now())

}
