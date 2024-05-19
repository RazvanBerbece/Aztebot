package messageEvent

import (
	"fmt"
	"strings"
	"time"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func UserRepReact(s *discordgo.Session, m *discordgo.MessageCreate) {

	phrase := m.Content // looks like "+rep @Usertag" / "-rep @Usertag"

	// exit this func early if not a rep action
	if !strings.Contains(phrase, "+rep") && !strings.Contains(phrase, "-rep") {
		return
	}

	tokens := strings.Split(phrase, " ")
	if len(tokens) < 2 {
		// invalidly formatted user rep action
		return
	}

	repModeInput := tokens[0] // +rep, -rep

	// Get user ID from Discord mention tag i.e <@1234>
	targetUserTag := tokens[1]
	if !strings.Contains(targetUserTag, "<@") {
		// invalidly formatted user rep targed ID
		return
	}

	targetUserId := utils.GetDiscordIdFromMentionFormat(targetUserTag)

	// Delay +- reps so the action can't be spammed
	mRepDelay := 5 // in minutes
	if timestamp, exists := globalState.LastUserReps[targetUserId]; exists {
		durationSinceRep := time.Since(timestamp)
		if int(durationSinceRep.Minutes()) < mRepDelay {
			// ignore it
			return
		}
	}

	userRepEntryExists := globalRepositories.UserRepRepository.EntryExists(targetUserId)
	switch userRepEntryExists {
	case -1:
		// err
		fmt.Println("An error ocurred while checking for user rep entry in the DB")
		return
	case 0:
		if repModeInput == "+rep" {
			err := globalRepositories.UserRepRepository.AddNewEntry(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding new user rep entry in the DB for %s: %v\n", targetUserId, err)
				return
			}
			err = globalRepositories.UserRepRepository.AddRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding rep to user in the DB for %s: %v\n", targetUserId, err)
				return
			}
			globalState.LastUserReps[targetUserId] = time.Now()
		} else if repModeInput == "-rep" {
			err := globalRepositories.UserRepRepository.AddNewEntry(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding new user rep entry in the DB for %s: %v\n", targetUserId, err)
				return
			}
			err = globalRepositories.UserRepRepository.RemoveRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while removing rep from user in the DB for %s: %v\n", targetUserId, err)
				return
			}
			globalState.LastUserReps[targetUserId] = time.Now()
		}
	case 1:
		if repModeInput == "+rep" {
			err := globalRepositories.UserRepRepository.AddRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding rep to user in the DB for %s: %v\n", targetUserId, err)
				return
			}
			globalState.LastUserReps[targetUserId] = time.Now()
		} else if repModeInput == "-rep" {
			err := globalRepositories.UserRepRepository.RemoveRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while removing rep from user in the DB for %s: %v\n", targetUserId, err)
				return
			}
			globalState.LastUserReps[targetUserId] = time.Now()
		}
	default:
		fmt.Printf("Multiple rep entries in the DB for %s\n", targetUserId)
		return
	}

}
