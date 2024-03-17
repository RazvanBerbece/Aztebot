package server_channel

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func VoiceChannelHasConnectedMembers(s *discordgo.Session, guildId string, channelId string) (bool, error) {

	guild, err := s.State.Guild(guildId)
	if err != nil {
		fmt.Println("Error retrieving guild:", err)
		return false, err
	}

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == channelId {
			return true, nil
		}
	}

	return false, nil
}

func CreateVoiceChannelForCategory(s *discordgo.Session, guildId string, categoryId string, channelName string, private bool) (*discordgo.Channel, error) {

	if private {
		channel, err := s.GuildChannelCreateComplex(guildId, discordgo.GuildChannelCreateData{
			Name:      channelName,
			Type:      discordgo.ChannelTypeGuildVoice,
			ParentID:  categoryId,
			UserLimit: 2,
		})

		if err != nil {
			fmt.Printf("An error ocurred while creating a dynamic private voice channel: %v\n", err)
			return nil, err
		}

		return channel, nil
	}

	channel, err := s.GuildChannelCreateComplex(guildId, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: categoryId,
	})

	if err != nil {
		fmt.Printf("An error ocurred while creating a dynamic voice channel: %v\n", err)
		return nil, err
	}

	return channel, nil
}
