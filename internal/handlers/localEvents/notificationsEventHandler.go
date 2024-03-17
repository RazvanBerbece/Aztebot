package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/notifications"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/bwmarrin/discordgo"
)

func HandleNotificationEvents(s *discordgo.Session) {

	for notificationEvent := range globalMessaging.NotificationsChannel {
		switch notificationEvent.Type {
		case "EMBED_WITH_TITLE_AND_FIELDS":
			err := notifications.SendNotificationWithFieldsToTextChannel(
				s,
				notificationEvent.TargetChannelId,
				*notificationEvent.Title,
				notificationEvent.Fields,
				*notificationEvent.UseThumbnail)
			if err != nil {
				fmt.Printf("Failed to process NotificationEvent (%s): %v\n", notificationEvent.Type, err)
			}
		case "EMBED_PASSTHROUGH":
			err := notifications.SendEmbedToTextChannel(
				s,
				notificationEvent.TargetChannelId,
				*notificationEvent.Embed)
			if err != nil {
				fmt.Printf("Failed to process NotificationEvent (%s): %v\n", notificationEvent.Type, err)
			}
		case "EMBED_WITH_ACTION_ROW":
			approvalMessageId, err := notifications.SendNotificationWithActionRowToTextChannel(
				s,
				notificationEvent.TargetChannelId,
				*notificationEvent.Title,
				notificationEvent.Fields,
				*notificationEvent.ActionRow,
				*notificationEvent.UseThumbnail)
			if err != nil {
				fmt.Printf("An error ocurred while sending confession approval notification: %v\n", err)
				return
			}

			// Keep to-be-approved confessions in-memory in order to forward them after approval
			globalState.ConfessionsToApprove[*approvalMessageId] = *notificationEvent.TextData

		default:
			fmt.Printf("Notification of type %s is not currently supported.", notificationEvent.Type)
		}
	}

}
