package dm

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendHelpDmToUser(session *discordgo.Session, originalMessage *discordgo.InteractionCreate, userId string, content string) error {

	channel, err := session.UserChannelCreate(userId)
	if err != nil {
		fmt.Println("error creating channel: ", err)
		session.ChannelMessageSend(
			originalMessage.ChannelID,
			"Something went wrong while sending the DM",
		)
		return err
	}

	_, err = session.ChannelMessageSend(channel.ID, content)
	if err != nil {
		fmt.Println("error sending DM message: ", err)
		session.ChannelMessageSend(
			originalMessage.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
		return err
	}

	return nil

}
