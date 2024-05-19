package messageEvent

import (
	"fmt"
	"strings"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
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
	targetUserId := utils.GetDiscordIdFromMentionFormat(targetUserTag)

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
		}
	case 1:
		if repModeInput == "+rep" {
			err := globalRepositories.UserRepRepository.AddRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding rep to user in the DB for %s: %v\n", targetUserId, err)
				return
			}
		} else if repModeInput == "-rep" {
			err := globalRepositories.UserRepRepository.RemoveRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while removing rep from user in the DB for %s: %v\n", targetUserId, err)
				return
			}
		}
	default:
		fmt.Printf("Multiple rep entries in the DB for %s\n", targetUserId)
		return
	}

}
