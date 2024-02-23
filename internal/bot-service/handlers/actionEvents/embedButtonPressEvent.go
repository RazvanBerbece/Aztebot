package actionEvent

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	supportSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/support"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleMessageComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		handleEmbedButtonPressEventHandler(s, i)
	}
}

func handleEmbedButtonPressEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	eventCustomId := i.MessageComponentData().CustomID

	switch eventCustomId {
	case "approve_confession":
		handleApproveConfession(s, i)
	case "decline_confession":
		handleDeclineConfession(s, i)
	}
}

func handleApproveConfession(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalApprovalMessageId := i.Message.ID
	originalApprovalMessageChannelId := i.Message.ChannelID

	// Get original interaction if it can be found in the in-memory map
	confessionMessage, exists := globals.ConfessionsToApprove[originalApprovalMessageId]
	if !exists {
		utils.SendErrorEmbedResponse(s, i.Interaction, "This confession message could not be found in the internal records.")
		return
	} else {
		// Send confession to designated text channel
		if channel, channelExists := globals.NotificationChannels["notif-confess"]; channelExists {
			go supportSlashHandlers.SendApprovedConfessionNotification(s, channel.ChannelId, confessionMessage)
		}
		delete(globals.ConfessionsToApprove, originalApprovalMessageChannelId)
	}

	// Respond to the button press
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You approved confession with ID `%s`. The message will now be forwarded to the designated channel", originalApprovalMessageId),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v\n", err)
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
	}

	// Cleanup
	go utils.DeleteInteractionResponse(s, i.Interaction, 3)
	go disableButtonsForApprovalActionRow(s, originalApprovalMessageChannelId, originalApprovalMessageId, "approve_confession", "decline_confession")
	go updateApprovedActionRowOriginalMessage(s, i.Member.User.Username, "APPROVED", originalApprovalMessageChannelId, originalApprovalMessageId, "approve_confession", "decline_confession")
}

func handleDeclineConfession(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalApprovalMessageId := i.Message.ID
	originalApprovalMessageChannelId := i.Message.ChannelID

	// Get original interaction if it can be found in the in-memory map
	_, exists := globals.ConfessionsToApprove[originalApprovalMessageId]
	if !exists {
		utils.SendErrorEmbedResponse(s, i.Interaction, "This confession message could not be found in the internal records.")
		return
	} else {
		delete(globals.ConfessionsToApprove, originalApprovalMessageId)
	}

	// Respond to the button press
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You declined confession with ID `%s`", originalApprovalMessageId),
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v\n", err)
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
	}

	// Cleanup
	go utils.DeleteInteractionResponse(s, i.Interaction, 3)
	go disableButtonsForApprovalActionRow(s, originalApprovalMessageChannelId, originalApprovalMessageId, "approve_confession", "decline_confession")
	go updateApprovedActionRowOriginalMessage(s, i.Member.User.Username, "DECLINED", originalApprovalMessageChannelId, originalApprovalMessageId, "approve_confession", "decline_confession")
}

func disableButtonsForApprovalActionRow(s *discordgo.Session, channelId string, messageId string, customApprovalEventId string, customDeclineEventId string) error {

	interactionEdit := &discordgo.MessageEdit{
		Channel: channelId,
		ID:      messageId,
		Content: nil,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "👍🏽",
						},
						Label:    "Accept",
						Style:    discordgo.SuccessButton,
						CustomID: customApprovalEventId,
						Disabled: true,
					},
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "👎🏽",
						},
						Label:    "Decline",
						Style:    discordgo.DangerButton,
						CustomID: customDeclineEventId,
						Disabled: true,
					},
				},
			},
		},
	}

	// Edit the original message with the disabled buttons in the action row
	_, err := s.ChannelMessageEditComplex(interactionEdit)
	if err != nil {
		// Handle error
		return err
	}

	return nil
}

func updateApprovedActionRowOriginalMessage(s *discordgo.Session, ownerTag string, opCode string, channelId string, messageId string, customApprovalEventId string, customDeclineEventId string) error {

	// Retrieve original embed
	message, err := s.ChannelMessage(channelId, messageId)
	if err != nil {
		return err
	}

	if len(message.Embeds) > 0 {

		originalEmbed := message.Embeds[0] // this gets mutated

		originalEmbedText := message.Embeds[0].Fields[0].Value
		updatedEmbedText := originalEmbedText + fmt.Sprintf("\n\n_(`%s` by `%s` at `%s`)_", opCode, ownerTag, time.Now())
		originalEmbed.Fields[0].Value = updatedEmbedText

		interactionEdit := &discordgo.MessageEdit{
			Channel: channelId,
			ID:      messageId,
			Content: nil,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "👍🏽",
							},
							Label:    "Accept",
							Style:    discordgo.SuccessButton,
							CustomID: customApprovalEventId,
							Disabled: true,
						},
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "👎🏽",
							},
							Label:    "Decline",
							Style:    discordgo.DangerButton,
							CustomID: customDeclineEventId,
							Disabled: true,
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{originalEmbed},
		}

		// Edit the original message with the disabled buttons in the action row
		_, err = s.ChannelMessageEditComplex(interactionEdit)
		if err != nil {
			// Handle error
			return err
		}
	}

	return nil
}
