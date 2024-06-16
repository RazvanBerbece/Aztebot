package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/domain"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
)

func ClearOldUserRepDelays() {

	var numSec int = 60 * 5           // run every 5 minutes
	threshold := time.Second * 60 * 5 // entries which are older than 5 minutes

	fmt.Println("[CRON] Starting Task ClearOldUserRepDelays() at", time.Now(), "running every", numSec, "seconds")

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go cleanupOldUserRepDelays(threshold)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func cleanupOldUserRepDelays(threshold time.Duration) {
	for repAuthor, reps := range globalState.LastGivenUserReps {
		// If an unused entry because all delays expired
		if len(reps) == 0 {
			// then remove the author from the map
			delete(globalState.LastGivenUserReps, repAuthor)
			continue
		}
		// Remove expired delays from the list of targets
		remainingTargets := []domain.GivenRep{}
		for _, rep := range reps {
			// If old enough
			if time.Since(time.Unix(int64(rep.Timestamp.Unix()), 0)) > threshold {
				continue
			}
			remainingTargets = append(remainingTargets, rep)
		}
		globalState.LastGivenUserReps[repAuthor] = remainingTargets
	}
}
