package messageEvent

import (
	"fmt"
	"time"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Increase stats for user
	err := globalsRepo.UserStatsRepository.IncrementMessagesSentForUser(m.Author.ID)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v", m.Author.ID, err)
	}

	err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(m.Author.ID)
	if err != nil {
		fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", m.Author.ID, err)
	}
	err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(m.Author.ID, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", m.Author.ID, err)
	}

}
