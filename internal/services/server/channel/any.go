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
		if channel.ParentID == categoryId {
			if strings.Contains(channel.Name, "~Extra~") {
				count++
			}
		}
	}

	return count, nil

}
