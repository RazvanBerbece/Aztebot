package messageEvent

import (
	"fmt"
	"strings"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/domain"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func UserRepReact(s *discordgo.Session, m *discordgo.MessageCreate) {

	authorUserId := m.Author.ID

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
	if repModeInput != "+rep" && repModeInput != "-rep" {
		// invalidly formatted user rep action
		return
	}

	// Get user ID from Discord mention tag i.e <@1234>
	targetUserTag := tokens[1]
	if !strings.Contains(targetUserTag, "<@") {
		// invalidly formatted user rep targed ID
		return
	}

	targetUserId := utils.GetDiscordIdFromMentionFormat(targetUserTag)

	// Don't allow users to =- rep themselves
	if targetUserId == m.Author.ID {
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		s.MessageReactionAdd(m.ChannelID, m.ID, "👺")
		return
	}

	// Delay +- reps so the action can't be spammed
	if targetReps, exists := globalState.LastGivenUserReps[authorUserId]; exists {
		if domain.IdInAuthorRepTargetList(targetUserId, targetReps) {
			// rep author already repped this target user
			s.MessageReactionAdd(m.ChannelID, m.ID, "⏳")
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

			// Register given rep to block follow-up reps to the same target from the same author
			reps := globalState.LastGivenUserReps[authorUserId]
			reps = append(reps, domain.GivenRep{
				To:        targetUserId,
				Timestamp: time.Now(),
			})
			globalState.LastGivenUserReps[authorUserId] = reps
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

			// Register given rep to block follow-up reps to the same target from the same author
			reps := globalState.LastGivenUserReps[authorUserId]
			reps = append(reps, domain.GivenRep{
				To:        targetUserId,
				Timestamp: time.Now(),
			})
			globalState.LastGivenUserReps[authorUserId] = reps
		}

		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

		// Reply with current rep status for target user
		rep, err := member.GetRep(targetUserId)
		if err != nil {
			fmt.Println("Failed to reply to user rep event")
			return
		}
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>: `%d` Rep", targetUserId, rep), m.Reference())
	case 1:
		if repModeInput == "+rep" {
			err := globalRepositories.UserRepRepository.AddRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while adding rep to user in the DB for %s: %v\n", targetUserId, err)
				return
			}

			// Register given rep to block follow-up reps to the same target from the same author
			reps := globalState.LastGivenUserReps[authorUserId]
			reps = append(reps, domain.GivenRep{
				To:        targetUserId,
				Timestamp: time.Now(),
			})
			globalState.LastGivenUserReps[authorUserId] = reps
		} else if repModeInput == "-rep" {
			err := globalRepositories.UserRepRepository.RemoveRep(targetUserId)
			if err != nil {
				fmt.Printf("An error ocurred while removing rep from user in the DB for %s: %v\n", targetUserId, err)
				return
			}

			// Register given rep to block follow-up reps to the same target from the same author
			reps := globalState.LastGivenUserReps[authorUserId]
			reps = append(reps, domain.GivenRep{
				To:        targetUserId,
				Timestamp: time.Now(),
			})
			globalState.LastGivenUserReps[authorUserId] = reps
		}

		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

		// Reply with current rep status for target user
		rep, err := member.GetRep(targetUserId)
		if err != nil {
			fmt.Println("Failed to reply to user rep event")
			return
		}
		_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>: `%d` Rep", targetUserId, rep), m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("Multiple rep entries in the DB for %s\n", targetUserId)
		return
	}

}
