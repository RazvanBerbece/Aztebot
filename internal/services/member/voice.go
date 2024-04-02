package member

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func DisconnectFromVoiceChannel(s *discordgo.Session, guildId string, userId string) error {
	err := s.GuildMemberMove(guildId, userId, nil)
	if err != nil {
		fmt.Printf("Failed to disconnect member from VC: %v\n", err)
		return err
	}
	return nil
}
