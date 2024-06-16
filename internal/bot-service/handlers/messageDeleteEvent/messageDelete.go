package messageEvent

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {

	deletedMessage := m.BeforeDelete

	if deletedMessage != nil {
		// Decrease stats for user
		err := globalsRepo.UserStatsRepository.DecrementMessagesSentForUser(deletedMessage.Author.ID)
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) message count: %v", deletedMessage.Author.ID, err)
		}
	}

}

// TODO: Add MessageDeleteBulk handler ?
