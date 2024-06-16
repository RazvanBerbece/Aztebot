package memberUpdateEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func MemberRoleUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	// DEBUG
	if globalConfiguration.AuditRoleUpdatesInChannel {
		currentRoles := m.Roles
		currentRolesString := ""
		for idx, roleId := range currentRoles {
			role, err := member.GetDiscordRole(s, m.GuildID, roleId)
			if err != nil {
				fmt.Printf("Error ocurred while retrieving Discord role: %v\n", err)
			}
			if idx < len(currentRoles)-1 {
				currentRolesString += fmt.Sprintf("%s,", role.Name)
			} else if idx == len(currentRoles)-1 {
				currentRolesString += fmt.Sprintf(role.Name)
			}
		}

		// Audit update by logging in provided debug channel
		if channel, channelExists := globalConfiguration.NotificationChannels["notif-debug"]; channelExists {
			content := fmt.Sprintf("Handling role update for %s [%s] (updated roles: %s)", m.Member.User.Username, m.Member.User.ID, currentRolesString)
			globalMessaging.NotificationsChannel <- events.NotificationEvent{
				TargetChannelId: channel.ChannelId,
				Type:            "DEFAULT",
				TextData:        &content,
			}
		}
	}

	// Sync user in DB with the current Discord member state
	err := member.SyncMember(s, m.GuildID, m.Member.User.ID, m.Member, globalConfiguration.OrderRoleNames, globalConfiguration.SyncProgressionInMemberUpdates)
	if err != nil {
		fmt.Printf("Error ocurred while syncing new user roles with the DB: %v\n", err)
	}
}
