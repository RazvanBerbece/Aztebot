package member

import (
	"time"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/bwmarrin/discordgo"
)

func DeleteMostRecentMemberMessages(s *discordgo.Session, guildId string, userId string, searchCount int, msgTimeLimit time.Duration) error {

	iterations := searchCount

	channels, err := s.GuildChannels(guildId)
	if err != nil {
		return err
	}

	for _, channel := range channels {

		// Don't delete messages from these channels (as they might be important for audits)
		if _, isDeletionExceptedChannel := globalConfiguration.DeleteExceptedChannels[channel.ID]; isDeletionExceptedChannel {
			continue
		}

		lastMessageId := "" // last retrieved message ID / also oldest message in the current batch

		// For multiple iterations of searches on the given channel
		for range iterations {

			// Retrieve a batch of messages from the given channel
			channelMessages, err := s.ChannelMessages(channel.ID, 100, lastMessageId, "", "")
			if err != nil {
				return err
			}

			if len(channelMessages) == 0 {
				// Empty channel
				continue
			}

			lastMessageId = channelMessages[len(channelMessages)-1].ID

			// Delete any messages belonging to the target user from the current batch
			for _, message := range channelMessages {
				if message.Author.ID == userId && time.Since(message.Timestamp) <= msgTimeLimit {
					err = s.ChannelMessageDelete(channel.ID, message.ID)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
