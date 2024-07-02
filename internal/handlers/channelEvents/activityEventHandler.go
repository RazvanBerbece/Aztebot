package channelHandlers

import (
	"fmt"
	"time"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

func HandleActivityRegistrationEvents() {

	for event := range globalMessaging.ActivityRegistrationsChannel {

		err := globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(event.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v\n", event.UserId, err)
		}
		err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(event.UserId, time.Now().Unix())
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v\n", event.UserId, err)
		}

		if event.Type != nil {
			switch *event.Type {
			case "MSG":
				err = globalRepositories.UserStatsRepository.IncrementMessagesSentForUser(event.UserId)
				if err != nil {
					fmt.Printf("An error ocurred while incrementing user (%s) message count: %v\n", event.UserId, err)
				}
			case "REACT":
				// nowt
			case "VC":
				// nowt
			case "MUSIC":
				// nowt
			case "SLASH":
				err := globalRepositories.UserStatsRepository.IncrementSlashCommandsUsedForUser(event.UserId)
				if err != nil {
					fmt.Printf("An error ocurred while incrementing user (%s) slash commands: %v\n", event.UserId, err)
				}
			default:
				fmt.Printf("Activity type %s not supported", *event.Type)
			}
		}
	}

}
