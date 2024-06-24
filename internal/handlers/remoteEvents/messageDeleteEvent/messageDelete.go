package messageEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
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
	authorIsBot, err := member.IsBot(s, globalConfiguration.DiscordMainGuildId, deletedMessageAuthorId, false)
	if err != nil {
		return
	}
	if authorIsBot == nil {
		return
	}
	if *authorIsBot {
		return
	}

	// Decrease stats for user
	err = globalRepositories.UserStatsRepository.DecrementMessagesSentForUser(deletedMessageAuthorId)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v\n", deletedMessageAuthorId, err)
	}

	// Remove experience points
	currentXp, err := member.RemoveMemberExperience(deletedMessageAuthorId, "MSG_REWARD")
	if err != nil {
		fmt.Printf("An error ocurred while removing message rewards (%d) from user (%s): %v\n", currentXp, deletedMessageAuthorId, err)
	}

	// Remove coins
	err = globalRepositories.WalletsRepository.SubtractFundsFromWallet(deletedMessageAuthorId, globalConfiguration.CoinReward_MessageSent)
	if err != nil {
		fmt.Printf("An error ocurred while removing message reward coins (%d) from user (%s): %v\n", currentXp, deletedMessageAuthorId, err)
	}

}

// TODO: Add MessageDeleteBulk handler ?
