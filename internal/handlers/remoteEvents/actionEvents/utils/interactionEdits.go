package actionEventsUtils

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func DisableButtonsForApprovalActionRow(s *discordgo.Session, channelId string, messageId string, customApprovalEventId string, customDeclineEventId string) error {

	interactionEdit := &discordgo.MessageEdit{
		Channel: channelId,
		ID:      messageId,
		Content: nil,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "👍🏽",
						},
						Label:    "Accept",
						Style:    discordgo.SuccessButton,
						CustomID: customApprovalEventId,
						Disabled: true,
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
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

func UpdateApprovedActionRowOriginalMessage(s *discordgo.Session, ownerTag string, opCode string, channelId string, messageId string, customApprovalEventId string, customDeclineEventId string) error {

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
			Components: &[]discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "👍🏽",
							},
							Label:    "Accept",
							Style:    discordgo.SuccessButton,
							CustomID: customApprovalEventId,
							Disabled: true,
						},
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
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
			Embeds: &[]*discordgo.MessageEmbed{originalEmbed},
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
