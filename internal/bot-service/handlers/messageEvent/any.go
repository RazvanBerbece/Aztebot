package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Increase stats for user
	userStatsRepo := repositories.NewUsersStatsRepository()
	err := userStatsRepo.IncrementMessagesSentForUser(m.Author.ID)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v", m.Author.ID, err)
	}

}
