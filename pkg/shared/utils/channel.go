package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ChannelHasDefaultInformationMessage(session *discordgo.Session, channelId string) (bool, error) {

	// Fetch message history for the specified channel
	messages, err := session.ChannelMessages(channelId, 10, "", "", "")
	if err != nil {
		fmt.Println("Error fetching messages: ", err)
		return false, err
	}

	if len(messages) > 0 && messages[0].Author.Bot {
		return true, nil
	} else {
		return false, nil
	}

}
