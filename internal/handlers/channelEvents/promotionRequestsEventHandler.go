package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func HandlePromotionRequestEvents(s *discordgo.Session, defaultOrderRoleNames []string, audit bool) {

	for xpEvent := range globalMessaging.PromotionRequestsChannel {

		var userGuildId = xpEvent.GuildId
		var userId = xpEvent.UserId
		var userTag = xpEvent.UserTag
		var userXp = xpEvent.CurrentXp
		var userCurrentLevel = xpEvent.CurrentLevel
		var userNumberMessagesSent = xpEvent.MessagesSent
		var userTimeSpentInVc = xpEvent.TimeSpentInVc

		const sHour = 60 * 60

		currentOrderRole, err := member.GetMemberOrderRole(userId, defaultOrderRoleNames)
		if err != nil {
			fmt.Printf("Error occurred while reading member order role from DB: %v\n", err)
		}

		// Check current stats against progression table
		// Figure out the promoted role to be given
		var promotedLevel int = 0
		var promotedRoleName string = ""
		switch {
		// No order
		case userXp < 7500:
			// no promotion required
			continue
		// First order
		case userXp >= 7500 && userXp < 10000:
			if userNumberMessagesSent >= 1000 && userTimeSpentInVc >= sHour*15 {
				promotedLevel = 1
				promotedRoleName = "🔗 Zelator"
			}
		case userXp >= 10000 && userXp < 15000:
			if userNumberMessagesSent >= 2500 && userTimeSpentInVc >= sHour*20 {
				promotedLevel = 2
				promotedRoleName = "📖 Theoricus"
			}
		case userXp >= 15000 && userXp < 30000:
			if userNumberMessagesSent >= 5000 && userTimeSpentInVc >= sHour*30 {
				promotedLevel = 3
				promotedRoleName = "📿 Philosophus"
			}
		// Second order
		case userXp >= 30000 && userXp < 45000:
			if userNumberMessagesSent >= 12500 && userTimeSpentInVc >= sHour*40 {
				promotedLevel = 4
				promotedRoleName = "🔮 Adeptus Minor"
			}
		case userXp >= 45000 && userXp < 50000:
			if userNumberMessagesSent >= 15000 && userTimeSpentInVc >= sHour*45 {
				promotedLevel = 5
				promotedRoleName = "〽️ Adeptus Major"
			}
		case userXp >= 50000 && userXp < 100000:
			if userNumberMessagesSent >= 20000 && userTimeSpentInVc >= sHour*50 {
				promotedLevel = 6
				promotedRoleName = "🧿 Adeptus Exemptus"
			}
		// Third order
		case userXp >= 100000 && userXp < 150000:
			if userNumberMessagesSent >= 35000 && userTimeSpentInVc >= sHour*200 {
				promotedLevel = 7
				promotedRoleName = "☀️ Magister Templi"
			}
		case userXp >= 150000 && userXp < 200000:
			if userNumberMessagesSent >= 45000 && userTimeSpentInVc >= sHour*250 {
				promotedLevel = 8
				promotedRoleName = "🧙🏼 Magus"
			}
		case userXp >= 200000:
			if userNumberMessagesSent >= 50000 && userTimeSpentInVc >= sHour*300 {
				promotedLevel = 9
				promotedRoleName = "⚔️ Ipsissimus"
			}
		}

		// This mismatch resolution is a result of the fact that progression roles were given outside the rules of progression
		// and now the bot has to resolve these mismatches.
		// Eventually, this code can and should be be removed.
		if currentOrderRole != nil {
			// mismatch between deserved role and actual role, so refresh
			if currentOrderRole.DisplayName != promotedRoleName && promotedRoleName != "" {
				// Give promoted level in DB
				err := globalRepositories.UsersRepository.SetLevel(userId, promotedLevel)
				if err != nil {
					fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
					continue
				}

				promotedRole, err := globalRepositories.RolesRepository.GetRole(promotedRoleName) // to append
				if err != nil {
					fmt.Printf("Error occurred while reading role %s from DB: %v\n", promotedRoleName, err)
					continue
				}

				err = member.RemoveAllMemberOrderRoles(userId, defaultOrderRoleNames)
				if err != nil {
					fmt.Printf("Error occurred while removing order member roles from DB: %v\n", err)
				}

				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
				}

				user, err := globalRepositories.UsersRepository.GetUser(userId)
				if err != nil {
					fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
				}
				err = member.RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
				if err != nil {
					fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
				}

				fmt.Printf("Resolved progression mismatch for %s ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, promotedLevel, promotedRole.DisplayName)

				continue
			}
		}

		// Promotion is available for current member
		if promotedLevel > userCurrentLevel {

			// Give promoted level in DB
			err := globalRepositories.UsersRepository.SetLevel(userId, promotedLevel)
			if err != nil {
				fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
				continue // skip event to allow retry with correct params
			}

			promotedRole, err := globalRepositories.RolesRepository.GetRole(promotedRoleName) // to append
			if err != nil {
				fmt.Printf("Error occurred while reading role from DB: %v\n", err)
				continue // skip event to allow retry with correct params
			}

			if promotedLevel == 1 {
				// no previous order role so no need to remove it, only append to list of IDs
				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
				}
			} else if promotedLevel > 1 {
				currentOrderRole, err := member.GetMemberOrderRole(userId, defaultOrderRoleNames) // to remove
				if err != nil {
					fmt.Printf("Error occurred while reading member order role from DB: %v\n", err)
				}
				if currentOrderRole != nil {
					// Only remove the current order role if one exists
					err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, currentOrderRole.Id)
					if err != nil {
						fmt.Printf("Error occurred while removing member role from DB: %v\n", err)
					}
				}
				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
				}
			}

			// Get refreshed role IDs after processing
			user, err := globalRepositories.UsersRepository.GetUser(userId)
			if err != nil {
				fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
			}
			err = member.RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
			if err != nil {
				fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
			}

			fmt.Printf("%s leveled up ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, promotedLevel, promotedRole.DisplayName)

			// Send notification and DM to audit progression
			if audit {
				go auditProgression(userId, promotedRole.DisplayName)
				go announceLevelUp(userId, promotedLevel, promotedRole.DisplayName)
			}
		}

	}

}

func auditProgression(userId string, newRoleName string) {
	// Audit by sending notification on designated channel
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-aztebotUpdatesChannel"]; channelExists {
		content := fmt.Sprintf("<@%s> has leveled up ! They attained the `%s` order role.", userId, newRoleName)
		globalMessaging.NotificationsChannel <- events.NotificationEvent{
			TargetChannelId: channel.ChannelId,
			Type:            "DEFAULT",
			TextData:        &content,
		}
	}
}

func announceLevelUp(userId string, newLevel int, newRoleName string) {
	dmEmbed := embed.NewEmbed().
		SetTitle("🤖⭐    Level up!").
		AddField("", fmt.Sprintf("You have officially attained the required activity metrics to progress to the `%s` order role. You are now level `%d`.", newRoleName, newLevel), false)

	globalMessaging.DirectMessagesChannel <- events.DirectMessageEvent{
		UserId: userId,
		Embed:  dmEmbed,
	}
}
