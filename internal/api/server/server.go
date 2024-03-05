package server

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/bwmarrin/discordgo"
)

func SetGlobalRestrictionsForRole(s *discordgo.Session, guildId string, targetRoleName string, exceptionChannelId *string) error {

	roleId := member.GetDiscordRoleIdForRoleWithName(s, guildId, targetRoleName)
	if roleId == nil {
		return fmt.Errorf("couldn't retrieve the Discord role ID for the Jailed feature")
	}

	channels, err := s.GuildChannels(guildId)
	if err != nil {
		return err
	}

	// Deny access to all channels and categories
	for _, channel := range channels {
		err := s.ChannelPermissionSet(channel.ID, *roleId, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionAll)
		if err != nil {
			fmt.Printf("Error setting denial permissions for channel %s: %v\n", channel.Name, err)
		}
	}

	// Allow basic perms for the exception channels - if applicable
	if exceptionChannelId != nil {
		err := s.ChannelPermissionSet(
			*exceptionChannelId,
			*roleId,
			discordgo.PermissionOverwriteTypeRole,
			discordgo.PermissionViewChannel|discordgo.PermissionSendMessages|discordgo.PermissionReadMessageHistory,
			0)
		if err != nil {
			fmt.Printf("Error setting denial permissions for channel with ID %s: %v\n", *exceptionChannelId, err)
		}
	}

	return nil
}
