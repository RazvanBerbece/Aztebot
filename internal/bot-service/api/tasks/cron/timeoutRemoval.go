package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func ClearExpiredTimeouts(s *discordgo.Session) {

	var numSec int
	if globals.TimeoutClearFrequencyErr != nil {
		numSec = 300
	} else {
		numSec = globals.TimeoutClearFrequency
	}

	fmt.Println("[CRON]\t\t\tStarting Task ClearExpiredTimeouts() at", time.Now(), "running every", numSec, "seconds")

	cleanupExpiredTimeouts() // initial run can happen at startup

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go cleanupExpiredTimeouts()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func cleanupExpiredTimeouts() {

	serverTimeouts, err := globalsRepo.TimeoutsRepository.GetAllTimeouts()
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
			err := globalsRepo.TimeoutsRepository.ArchiveTimeout(timeout.UserId, timeout.Reason, expiryTime.Unix())
			if err != nil {
				fmt.Printf("An error ocurred while archiving warning for user %s: %v", timeout.UserId, err)
			}
			// Remove it from the active timeouts table
			err = globalsRepo.TimeoutsRepository.ClearTimeoutForUser(timeout.UserId)
			if err != nil {
				fmt.Printf("An error ocurred while clearing an expired warning from user %s: %v", timeout.UserId, err)
			}
		}

	}

}
