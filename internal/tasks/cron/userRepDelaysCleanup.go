package cron

import (
	"fmt"
	"time"

	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
)

func ClearOldUserRepDelays() {

	var numSec int = 60 * 10           // run every 10 minutes
	threshold := time.Second * 60 * 10 // entries which are older than 10 minutes

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
	for userId, timestamp := range globalState.LastUserReps {
		// If old enough
		if time.Since(time.Unix(int64(timestamp.Unix()), 0)) > threshold {
			delete(globalState.LastUserReps, userId)
		}
	}
}
