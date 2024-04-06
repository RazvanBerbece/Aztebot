package server_channel

import (
	"strings"

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

func GetNumberOfDynamicChannelsForCategory(s *discordgo.Session, guildId string, categoryId string) (int, error) {

	count := 0

	channels, err := s.GuildChannels(guildId)
	if err != nil {
		return -1, err
	}

	for _, channel := range channels {
		// If in the provided category, a voice channel and a dynamic one
		if channel.ParentID == categoryId && channel.Type == discordgo.ChannelTypeGuildVoice && strings.Contains(channel.Name, "~Extra~") {
			// then count it
			count++
		}
	}

	return count, nil

}
