package guildRemoveEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func GuildRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	if globalConfiguration.AuditMemberDeletesInChannel {
		if channel, channelExists := globalConfiguration.NotificationChannels["notif-debug"]; channelExists {
			content := fmt.Sprintf("%s left the server", m.Member.User.Username)
			globalMessaging.NotificationsChannel <- events.NotificationEvent{
				TargetChannelId: channel.ChannelId,
				Type:            "DEFAULT",
				TextData:        &content,
			}
		}
	}

	// Delete user from all tables
	userId := m.Member.User.ID
	err := member.DeleteAllMemberData(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s data from DB tables on kick action: %v", userId, err)
		return
	}

	// Other actions to do on guild leave

}
