package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {

	deletedMessage := m.BeforeDelete
	if deletedMessage == nil {
		return
	}
	if deletedMessage.Author == nil {
		return
	}
	deletedMessageAuthorId := deletedMessage.Author.ID

	// Ignore all messages created by bots
	authorIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, deletedMessageAuthorId, false)
	if err != nil {
		return
	}
	if authorIsBot == nil {
		return
	}
	if *authorIsBot {
		return
	}

	if deletedMessage != nil {
		// Decrease stats for user
		err := globalsRepo.UserStatsRepository.DecrementMessagesSentForUser(deletedMessageAuthorId)
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) message count: %v", deletedMessageAuthorId, err)
		}

		// Remove experience points
		currentXp, err := member.RemoveMemberExperience(deletedMessageAuthorId, "MSG_REWARD")
		if err != nil {
			fmt.Printf("An error ocurred while removing message rewards (%d) from user (%s): %v", currentXp, deletedMessageAuthorId, err)
		}
	}

}

// TODO: Add MessageDeleteBulk handler ?
