package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
)

func ProcessRemoveArchivedTimeouts(months int) {

	// TODO: Move this into env variable and use cron syntax
	var numSec int = 2.628e+6 // run every month

	fmt.Println("[CRON] Starting Cron Ticker ProcessRemoveArchivedTimeouts() at", time.Now(), "running every month")

	// Inject new connections
	timeoutsRepository := repositories.NewTimeoutsRepository()

	go CleanupArchivedTimeouts(timeoutsRepository)

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go CleanupArchivedTimeouts(timeoutsRepository)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func CleanupArchivedTimeouts(timeoutsRepository *repositories.TimeoutsRepository) {

	fmt.Println("[CRON] Starting Task CleanupArchivedTimeouts() at", time.Now())

	archivedTimeouts, err := timeoutsRepository.GetAllArchivedTimeouts()
	if err != nil {
		fmt.Println("[CRON] Failed Task CleanupArchivedTimeouts() at", time.Now(), "with error", err)
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
				fmt.Println("[CRON] Failed Task CleanupArchivedTimeouts() at", time.Now(), "with error", err)
			}
		}
	}

	fmt.Println("[CRON] Finished Task CleanupArchivedTimeouts() at", time.Now())

}
