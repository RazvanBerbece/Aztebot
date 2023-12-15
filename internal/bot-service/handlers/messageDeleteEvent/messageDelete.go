package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/bwmarrin/discordgo"
)

func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {

	deletedMessage := m.BeforeDelete

	if deletedMessage != nil {
		// Decrease stats for user
		userStatsRepo := repositories.NewUsersStatsRepository()
		err := userStatsRepo.DecrementMessagesSentForUser(deletedMessage.Author.ID)
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) message count: %v", deletedMessage.Author.ID, err)
		}
	}

}

// TODO: Add MessageDeleteBulk handler ?
