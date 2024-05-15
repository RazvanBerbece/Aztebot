package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandlePromotionRequestEvents(s *discordgo.Session) {

	for xpEvent := range globalMessaging.PromotionRequestsChannel {

		var userGuildId = xpEvent.GuildId
		var userId = xpEvent.UserId
		var userXp = xpEvent.CurrentXp
		var userNumberMessagesSent = xpEvent.MessagesSent
		var userTimeSpentInVc = xpEvent.TimeSpentInVc

		// Check current stats against progression table
		// Figure out the promoted role to be given
		var promotedRoleId = -1
		switch {
		// No order
		case userXp < 7500:
		// First order
		case userXp >= 7500 && userXp < 10000:
		case userXp >= 10000 && userXp < 15000:
		case userXp >= 15000 && userXp < 20000:
		case userXp >= 20000 && userXp < 30000:
		// Second order
		case userXp >= 30000 && userXp < 45000:
		case userXp >= 45000 && userXp < 50000:
		case userXp >= 50000 && userXp < 100000:
		// Third order
		case userXp >= 100000 && userXp < 150000:
		case userXp >= 150000 && userXp < 200000:
		case userXp >= 200000:
		}

		// Give promoted role to member

		// Refresh order role
		err := member.RefreshDiscordOrderRoleForMember(s, userGuildId, userId)
		if err != nil {
			fmt.Printf("Error ocurred while refreshing member order role on-Discord: %v\n", err)
		}

	}

}
