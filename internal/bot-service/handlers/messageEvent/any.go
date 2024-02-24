package messageEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageCreatorUserId := m.Author.ID

	// Ignore all messages created by the bot itself
	if messageCreatorUserId == s.State.User.ID {
		return
	}

	// Increase stats for user
	err := globalsRepo.UserStatsRepository.IncrementMessagesSentForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v\n", messageCreatorUserId, err)
	}

	err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v\n", messageCreatorUserId, err)
	}
	err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(messageCreatorUserId, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v\n", m.Author.ID, err)
	}

	// Grant experience points
	currentXp, err := member.GrantMemberExperience(messageCreatorUserId, "MSG_REWARD", nil)
	if err != nil {
		fmt.Printf("An error ocurred while granting message rewards (%d) to user (%s): %v\n", currentXp, messageCreatorUserId, err)
	}

}
