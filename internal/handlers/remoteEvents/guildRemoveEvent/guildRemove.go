package guildRemoveEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func GuildRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	if globalConfiguration.AuditMemberDeletesInChannel {
		logMsg := fmt.Sprintf("`%s` left the server", m.Member.User.Username)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		discordChannelLogger.LogInfo(logMsg)
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
