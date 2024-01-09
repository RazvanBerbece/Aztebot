package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ChannelHasDefaultInformationMessage(session *discordgo.Session, channelId string) (bool, error) {

	// Fetch message history for the specified channel
	messages, err := session.ChannelMessages(channelId, 10, "", "", "") // returns messages in descending order of the timestamp (new messages come first)
	if err != nil {
		fmt.Println("Error fetching messages: ", err)
		return false, err
	}

	if len(messages) > 0 && messages[len(messages)-1].Author.Bot {
		return true, nil
	} else {
		return false, nil
	}

}

func TargetChannelIsForMusicListening(musicChannels map[string]string, channelId string) bool {
	for id := range musicChannels {
		if channelId == id {
			// Target VC is a music-specific channel
			return true
		}
	}
	return false
}
