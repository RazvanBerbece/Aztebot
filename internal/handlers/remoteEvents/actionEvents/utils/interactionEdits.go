package actionEventsUtils

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/domain"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
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

func UpdatePaginatedEmbedPage(s *discordgo.Session, embedData *domain.EmbedData, opCode string, channelId string, messageId string, pageSize int) error {

	// Retrieve original embed (the one with the pagination action row)
	message, err := s.ChannelMessage(channelId, messageId)
	if err != nil {
		return err
	}

	if len(message.Embeds) > 0 {

		pages := (len(*embedData.FieldData) + pageSize - 1) / pageSize

		// Calculate next page (and wrap if necessary)
		currentPage := embedData.CurrentPage
		if opCode == "NEXT" {
			currentPage += 1
		} else if opCode == "PREV" {
			currentPage -= 1
		}
		if currentPage > pages {
			currentPage = 1
		} else if currentPage < 1 {
			currentPage = pages
		}

		// Update map to hold new page number
		// assume that key exists
		globalState.EmbedsToPaginate[messageId] = domain.EmbedData{
			ChannelId:   embedData.ChannelId,
			CurrentPage: currentPage,
			FieldData:   embedData.FieldData,
			Timestamp:   embedData.Timestamp,
		}

		originalEmbed := message.Embeds[0] // this gets mutated

		// Determine the start and end index of fields to display for the current page
		startIdx := (currentPage - 1) * pageSize
		endIdx := startIdx + pageSize
		if endIdx > len(*embedData.FieldData) {
			endIdx = len(*embedData.FieldData)
		}

		fields := *embedData.FieldData
		paginatedFields := fields[startIdx:endIdx]
		originalEmbed.Fields = paginatedFields
		originalEmbed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d / %d", currentPage, pages),
		}
		interactionEdit := &discordgo.MessageEdit{
			Channel: channelId,
			ID:      messageId,
			Content: nil,
			Embeds:  &[]*discordgo.MessageEmbed{originalEmbed},
			Components: &[]discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "⬅️",
							},
							Label:    "Previous",
							Style:    discordgo.PrimaryButton,
							CustomID: globalMessaging.PreviousPageOnEmbedEventId,
							Disabled: false,
						},
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "➡️",
							},
							Label:    "Next",
							Style:    discordgo.PrimaryButton,
							CustomID: globalMessaging.NextPageOnEmbedEventId,
							Disabled: false,
						},
					},
				},
			},
		}

		_, err = s.ChannelMessageEditComplex(interactionEdit)
		if err != nil {
			// Handle error
			return err
		}
	}

	return nil
}

func DisablePaginatedEmbed(s *discordgo.Session, channelId string, messageId string) error {

	// Retrieve original embed (the one with the pagination action row)
	message, err := s.ChannelMessage(channelId, messageId)
	if err != nil {
		return err
	}

	originalEmbed := message.Embeds[0]

	interactionEdit := &discordgo.MessageEdit{
		Channel: channelId,
		ID:      messageId,
		Content: nil,
		Embeds:  &[]*discordgo.MessageEmbed{originalEmbed},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "⬅️",
						},
						Label:    "Previous",
						Style:    discordgo.PrimaryButton,
						CustomID: globalMessaging.PreviousPageOnEmbedEventId,
						Disabled: true,
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "➡️",
						},
						Label:    "Next",
						Style:    discordgo.PrimaryButton,
						CustomID: globalMessaging.NextPageOnEmbedEventId,
						Disabled: true,
					},
				},
			},
		},
	}

	_, err = s.ChannelMessageEditComplex(interactionEdit)
	if err != nil {
		// Handle error
		return err
	}

	return nil
}
