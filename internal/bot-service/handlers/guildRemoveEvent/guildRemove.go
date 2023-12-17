package guildRemoveEvent

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func GuildRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	logging.LogHandlerCall("GuildRemove", "")

	// Delete user from all tables
	userId := m.Member.User.ID
	err := globalsRepo.UserStatsRepository.DeleteUserStats(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s stats from DB: %v", userId, err)
	}
	err = globalsRepo.UsersRepository.DeleteUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s from DB: %v", userId, err)
	}

	// Other actions to do on guild leave

}
