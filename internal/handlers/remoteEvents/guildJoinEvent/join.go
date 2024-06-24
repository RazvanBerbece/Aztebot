package joinEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
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
		logMsg := fmt.Sprintf("`%s` joined the OTA server", m.Member.User.Username)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogInfo(logMsg)
	}

	// Store newly-joined user to DB tables (probably only the initial details and awaiting for verification and cron sync)
	go member.SyncMember(s, globalConfiguration.DiscordMainGuildId, m.Member.User.ID, m.Member, globalConfiguration.OrderRoleNames, false)

	// Other actions to do on guild join
	// e.g guide DM, etc.

}
