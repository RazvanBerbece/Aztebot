package dm

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func DmUser(session *discordgo.Session, userId string, content string) error {

	channel, err := session.UserChannelCreate(userId)
	if err != nil {
		fmt.Println("error creating DM channel: ", err)
		return err
	}

	_, err = session.ChannelMessageSend(channel.ID, content)
	if err != nil {
		fmt.Println("error sending DM message: ", err)
		return err
	}

	return nil

}

func DmEmbedUser(session *discordgo.Session, userId string, embed discordgo.MessageEmbed) error {

	channel, err := session.UserChannelCreate(userId)
	if err != nil {
		fmt.Println("error creating embed DM channel: ", err)
		return err
	}

	_, err = session.ChannelMessageSendEmbed(channel.ID, &embed)
	if err != nil {
		fmt.Println("error sending embed DM message: ", err)
		return err
	}

	return nil

}
