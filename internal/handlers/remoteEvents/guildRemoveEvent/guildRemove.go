package guildRemoveEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func GuildRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	userId := m.Member.User.ID

	if globalConfiguration.AuditMemberDeletesInChannel {
		logMsg := fmt.Sprintf("`%s` [%s] left the server", m.Member.User.Username, userId)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogInfo(logMsg)
	}

	// Clear all stored member data
	err := member.DeleteAllMemberData(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s data from DB tables on kick action: %v", userId, err)
		return
	}

	// Clear all local member data
	delete(globalState.MusicSessions, userId)
	delete(globalState.VoiceSessions, userId)
	delete(globalState.StreamSessions, userId)
	delete(globalState.DeafSessions, userId)

}
