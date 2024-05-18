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

		// Check current stats against progression table
		// Figure out the promoted role to be given
		processedRoleName, processedLevel := member.GetRoleNameAndLevelFromStats(userXp, userNumberMessagesSent, userTimeSpentInVc)

		currentOrderRoles, err := member.GetMemberOrderRoles(userId, defaultOrderRoleNames)
		if err != nil {
			fmt.Printf("Error occurred while reading member order role from DB: %v\n", err)
			continue
		}

		if len(currentOrderRoles) > 1 {
			if processedRoleName != "" {
				// resolve multiple role mismatches
				err = member.RemoveAllMemberOrderRoles(userId, defaultOrderRoleNames)
				if err != nil {
					fmt.Printf("Error occurred while removing order member roles from DB: %v\n", err)
				}

				err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
				if err != nil {
					fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
					continue
				}

				processedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
				if err != nil {
					fmt.Printf("Error occurred while reading role %s from DB: %v\n", processedRoleName, err)
					continue
				}

				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, processedRole.Id)
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

				fmt.Printf("Resolved progression mismatch (0) for %s ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, processedLevel, processedRole.DisplayName)

				continue
			}
		} else if len(currentOrderRoles) == 1 {

			// This mismatch resolution is a result of the fact that progression roles were given outside the rules of progression
			// and now the bot has to resolve these mismatches.
			// Eventually, this code can and should be be removed.

			// mismatch between deserved role and actual role, so refresh
			if processedRoleName != "" && currentOrderRoles[0].DisplayName != processedRoleName {
				// Give promoted level in DB
				err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
				if err != nil {
					fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
					continue
				}

				promotedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
				if err != nil {
					fmt.Printf("Error occurred while reading role %s from DB: %v\n", processedRoleName, err)
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

				fmt.Printf("Resolved progression mismatch (1) for %s ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, processedLevel, promotedRole.DisplayName)

				continue
			}
		} else if len(currentOrderRoles) == 0 {
			if processedRoleName != "" {
				processedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
				if err != nil {
					fmt.Printf("Error occurred while reading role %s from DB: %v\n", processedRoleName, err)
					continue
				}

				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, processedRole.Id)
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

				fmt.Printf("Resolved progression mismatch (2) for %s ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, processedLevel, processedRole.DisplayName)

				continue
			}
		}

		// Promotion is available for current member (and no mismatch was detected)
		if processedLevel > userCurrentLevel && processedRoleName != "" {

			// Give promoted level in DB
			err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
			if err != nil {
				fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
				continue // skip event to allow retry with correct params
			}

			promotedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
			if err != nil {
				fmt.Printf("Error occurred while reading role from DB: %v\n", err)
				continue // skip event to allow retry with correct params
			}

			if processedLevel == 1 {
				// no previous order role so no need to remove it, only append to list of IDs
				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
				}
			} else if processedLevel > 1 {
				for _, orderRole := range currentOrderRoles {
					err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, orderRole.Id)
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

			fmt.Printf("%s leveled up ! (%d -> %d) | New role: %s\n", userTag, userCurrentLevel, processedLevel, promotedRole.DisplayName)

			// Send notification and DM to audit progression
			if audit {
				go auditProgression(userId, promotedRole.DisplayName)
				go announceLevelUp(userId, processedLevel, promotedRole.DisplayName)
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
