package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/notifications"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
)

func HandleNotificationEvents() {

	for notificationEvent := range globals.NotificationsChannel {
		switch notificationEvent.Type {
		case "EMBED_WITH_TITLE_AND_FIELDS":
			err := notifications.SendNotificationToTextChannel(
				notificationEvent.Session,
				notificationEvent.TargetChannelId,
				*notificationEvent.Title,
				notificationEvent.Fields,
				*notificationEvent.UseThumbnail)
			if err != nil {
				fmt.Printf("Failed to process NotificationEvent (%s): %v\n", notificationEvent.Type, err)
			}
		case "EMBED_PASSTHROUGH":
			err := notifications.SendEmbedToTextChannel(
				notificationEvent.Session,
				notificationEvent.TargetChannelId,
				*notificationEvent.Embed)
			if err != nil {
				fmt.Printf("Failed to process NotificationEvent (%s): %v\n", notificationEvent.Type, err)
			}
		case "EMBED_WITH_ACTION_ROW":
			approvalMessageId, err := notifications.SendNotificationWithActionRowToTextChannel(
				notificationEvent.Session,
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
			globals.ConfessionsToApprove[*approvalMessageId] = *notificationEvent.TextData

		default:
			fmt.Printf("Notification of type %s is not currently supported.", notificationEvent.Type)
		}
	}

}
