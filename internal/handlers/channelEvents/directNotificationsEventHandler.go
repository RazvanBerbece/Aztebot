package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/notifications"
	"github.com/bwmarrin/discordgo"
)

func HandleDirectMessageEvents(s *discordgo.Session) {

	for directMessageEvent := range globalMessaging.DirectMessagesChannel {
		switch directMessageEvent.Type {
		case "EMBED_WITH_TITLE_AND_FIELDS":
			err := notifications.SendNotificationWithFieldsToTextChannel(
				s,
				directMessageEvent.TargetChannelId,
				*directMessageEvent.Title,
				directMessageEvent.Fields,
				*directMessageEvent.UseThumbnail)
			if err != nil {
				fmt.Printf("Failed to process directMessageEvent (%s): %v\n", directMessageEvent.Type, err)
			}
		case "EMBED_PASSTHROUGH":
			err := notifications.SendEmbedToTextChannel(
				s,
				directMessageEvent.TargetChannelId,
				*directMessageEvent.Embed)
			if err != nil {
				fmt.Printf("Failed to process directMessageEvent (%s): %v\n", directMessageEvent.Type, err)
			}
		default:
			fmt.Printf("Notification of type %s is not currently supported.", directMessageEvent.Type)
		}
	}

}
