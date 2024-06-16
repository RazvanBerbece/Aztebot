package server_channel

import "github.com/bwmarrin/discordgo"

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
