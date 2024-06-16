package cron

import (
	"fmt"
	"time"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/bwmarrin/discordgo"
)

func ProcessClearExpiredTimeouts(s *discordgo.Session) {

	var numSec int
	if globalConfiguration.TimeoutClearFrequencyErr != nil {
		numSec = 300
	} else {
		numSec = globalConfiguration.TimeoutClearFrequency
	}

	fmt.Println("[CRON] Starting Task ClearExpiredTimeouts() at", time.Now(), "running every", numSec, "seconds")

	// Inject repositories
	timeoutsRepository := repositories.NewTimeoutsRepository()

	go cleanupExpiredTimeouts(timeoutsRepository) // initial run can happen at startup

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go cleanupExpiredTimeouts(timeoutsRepository)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func cleanupExpiredTimeouts(timeoutsRepository *repositories.TimeoutsRepository) {

	serverTimeouts, err := timeoutsRepository.GetAllTimeouts()
	if err != nil {
		fmt.Printf("An error ocurred while retrieving server timeouts: %v", err)
	}

	for _, timeout := range serverTimeouts {

		timeoutCreationTime := time.Unix(timeout.CreationTimestamp, 0)
		duration := time.Second * time.Duration(timeout.SDuration)
		expiryTime := timeoutCreationTime.Add(duration)

		// If the timeout is expired
		if expiryTime.Before(time.Now()) {
			// Archive it
			err := timeoutsRepository.ArchiveTimeout(timeout.UserId, timeout.Reason, expiryTime.Unix())
			if err != nil {
				fmt.Printf("An error ocurred while archiving warning for user %s: %v", timeout.UserId, err)
			}
			// Remove it from the active timeouts table and member
			err = timeoutsRepository.ClearTimeoutForUser(timeout.UserId)
			if err != nil {
				fmt.Printf("An error ocurred while clearing an expired warning from user %s: %v", timeout.UserId, err)
			}
		}

	}

}
