package supportSlashHandlers

import (
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashConfess(s *discordgo.Session, i *discordgo.InteractionCreate) {

	message := i.ApplicationCommandData().Options[0].StringValue()

	// Send notification to target channel to announce the confession
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-confessApproval"]; channelExists {
		SendConfessionApprovalNotification(s, channel.ChannelId, message)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SimpleEmbed("", "Confession submitted for approval."),
			},
		})

		utils.DeleteInteractionResponse(s, i.Interaction, 0)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SimpleEmbed("", "No confession channel was found."),
			},
		})

		utils.DeleteInteractionResponse(s, i.Interaction, 0)
	}
}

func SendConfessionApprovalNotification(s *discordgo.Session, channelId string, message string) {

	fields := []discordgo.MessageEmbedField{
		{
			Name:   "",
			Value:  message,
			Inline: false,
		},
	}

	// Add action row with approval/disproval buttons to the confession approval embed being posted
	actionRow := embed.GetApprovalActionRowForEmbed(globalMessaging.ConfessionApprovalEventId, globalMessaging.ConfessionDisprovalEventId)
	notificationTitle := "New `/confess` to Approve"
	useThumbnail := false
	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Title:           &notificationTitle,
		Type:            "EMBED_WITH_ACTION_ROW",
		Fields:          fields,
		ActionRow:       &actionRow,
		TextData:        &message,
		UseThumbnail:    &useThumbnail,
	}

}

func SendApprovedConfessionNotification(s *discordgo.Session, channelId string, message string) {

	fields := []discordgo.MessageEmbedField{
		{
			Name:   "",
			Value:  message,
			Inline: false,
		},
	}

	emptyTitle := ""
	useThumbnail := false
	authorName := "New Confession Published"
	authorAvatarUrl := "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg"
	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Title:           &emptyTitle,
		Type:            "EMBED_WITH_TITLE_AND_FIELDS",
		Fields:          fields,
		UseThumbnail:    &useThumbnail,
		AuthorName:      &authorName,
		AuthorAvatarURL: &authorAvatarUrl,
	}

}
