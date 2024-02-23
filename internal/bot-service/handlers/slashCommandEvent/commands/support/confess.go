package supportSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/notifications"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashConfess(s *discordgo.Session, i *discordgo.InteractionCreate) {

	message := i.ApplicationCommandData().Options[0].StringValue()

	// Send notification to target channel to announce the confession
	if channel, channelExists := globals.NotificationChannels["notif-confessApproval"]; channelExists {
		go SendConfessionApprovalNotification(s, channel.ChannelId, message)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("", "Confession submitted for approval."),
		},
	})
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
	actionRow := embed.GetApprovalActionRowForEmbed(globals.ConfessionApprovalEventId, globals.ConfessionDisprovalEventId)
	approvalMessageId, err := notifications.SendNotificationWithActionRowToTextChannel(s, channelId, "New `/confess` to Approve", fields, actionRow, false)
	if err != nil {
		fmt.Printf("An error ocurred while sending confession approval notification: %v\n", err)
		return
	}

	// Keep to-be-approved confessions in-memory in order to forward them after approval
	globals.ConfessionsToApprove[*approvalMessageId] = message

}

func SendApprovedConfessionNotification(s *discordgo.Session, channelId string, message string) {

	fields := []discordgo.MessageEmbedField{
		{
			Name:   "",
			Value:  message,
			Inline: false,
		},
	}

	notifications.SendNotificationToTextChannel(s, channelId, "", fields, false)

}
