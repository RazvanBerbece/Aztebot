package guildRemoveEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func GuildRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	utils.LogHandlerCall("GuildRemove", "")

	// Delete user from all tables
	userId := m.Member.User.ID
	err := member.DeleteAllMemberData(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s data from DB tables on kick action: %v", userId, err)
		return
	}

	// Other actions to do on guild leave

}
