package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {

	deletedMessage := m.BeforeDelete
	if deletedMessage.Author == nil {
		// Probably an embed, so ignore
		return
	}
	deletedMessageAuthor := deletedMessage.Author.ID

	if deletedMessage != nil {
		// Decrease stats for user
		err := globalsRepo.UserStatsRepository.DecrementMessagesSentForUser(deletedMessageAuthor)
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) message count: %v", deletedMessageAuthor, err)
		}

		// Remove experience points
		currentXp, err := member.RemoveMemberExperience(deletedMessageAuthor, "MSG_REWARD")
		if err != nil {
			fmt.Printf("An error ocurred while removing message rewards (%d) from user (%s): %v", currentXp, deletedMessageAuthor, err)
		}
	}

}

// TODO: Add MessageDeleteBulk handler ?
