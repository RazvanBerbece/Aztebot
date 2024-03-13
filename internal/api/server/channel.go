package server

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetCategoryIdForChannel(s *discordgo.Session, guildId string, channelId string) (string, error) {

	channels, err := s.GuildChannels(guildId)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.ID == channelId {
			// The ParentID should be the category - if it exists
			return channel.ParentID, nil
		}
	}

	return "", nil

}

func CreateVoiceChannelForCategory(s *discordgo.Session, guildId string, categoryId string, channelName string, private bool) (*discordgo.Channel, error) {

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