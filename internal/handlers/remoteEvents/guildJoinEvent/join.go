package joinEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a new joined member.
func GuildJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	// Audit member join by logging in provided debug channel
	if globalConfiguration.AuditMemberJoinsInChannel {
		if channel, channelExists := globalConfiguration.NotificationChannels["notif-debug"]; channelExists {
			content := fmt.Sprintf("%s joined the OTA server", m.Member.User.Username)
			globalMessaging.NotificationsChannel <- events.NotificationEvent{
				TargetChannelId: channel.ChannelId,
				Type:            "DEFAULT",
				TextData:        &content,
			}
		}
	}

	// Store newly-joined user to DB tables (probably only the initial details and awaiting for verification and cron sync)
	err := member.SyncMember(s, globalConfiguration.DiscordMainGuildId, m.Member.User.ID, m.Member, globalConfiguration.OrderRoleNames, globalConfiguration.SyncProgressionInMemberUpdates)
	if err != nil {
		fmt.Printf("Error storing new member %s to DB: %v", m.Member.User.Username, err)
	}

	// Other actions to do on guild join
	// e.g guide DM, etc.

}
