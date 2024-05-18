package notifications

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func SendDefaultEmbedToTextChannel(s *discordgo.Session, channelId string, content string) error {

	embed := embed.NewEmbed().
		SetColor(000000)

	embed.AddField("", content, false)

	_, err := s.ChannelMessageSendEmbed(channelId, embed.MessageEmbed)
	if err != nil {
		fmt.Printf("Error sending notification to channel %s: %v", channelId, err)
		return err
	}

	return nil

}

func SendEmbedToTextChannel(s *discordgo.Session, channelId string, embed embed.Embed) error {

	_, err := s.ChannelMessageSendEmbed(channelId, embed.MessageEmbed)
	if err != nil {
		fmt.Printf("Error sending embed to channel %s: %v", channelId, err)
		return err
	}

	return nil

}

func SendNotificationWithFieldsToTextChannel(s *discordgo.Session, channelId string, notificationTitle string, fields []discordgo.MessageEmbedField, useThumbnail bool, authorName *string, authorAvatarUrl *string) error {

	// Build notification embed
	embed := embed.NewEmbed().
		SetColor(000000)

	if authorName != nil && authorAvatarUrl != nil {
		embed.SetAuthor(*authorName, *authorAvatarUrl)
	}

	// Don't show feedback bot emojis when there is no title,
	// as usually a notification with no title is meant to be kept minimalistic
	if notificationTitle != "" {
		embed.SetTitle(fmt.Sprintf("🤖ℹ️   %s", notificationTitle))
	}

	if useThumbnail {
		embed.SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg")
	}

	for _, field := range fields {
		embed.AddField(field.Name, field.Value, field.Inline)
	}

	_, err := s.ChannelMessageSendEmbed(channelId, embed.MessageEmbed)
	if err != nil {
		fmt.Printf("Error sending notification to channel %s: %v", channelId, err)
		return err
	}

	return nil

}

func SendNotificationWithActionRowToTextChannel(s *discordgo.Session, channelId string, notificationTitle string, fields []discordgo.MessageEmbedField, actionsRow discordgo.ActionsRow, useThumbnail bool, authorName *string, authorAvatarUrl *string) (*string, error) {

	// Build notification embed
	embed := embed.NewEmbed().
		SetColor(000000)

	if authorName != nil && authorAvatarUrl != nil {
		embed.SetAuthor(*authorName, *authorAvatarUrl)
	}

	// Don't show feedback bot emojis when there is no title,
	// as usually a notification with n otitle is meant to be kept minimalistic
	if notificationTitle != "" {
		embed.SetTitle(fmt.Sprintf("🤖ℹ️   %s", notificationTitle))
	}

	if useThumbnail {
		embed.SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg")
	}

	for _, field := range fields {
		embed.AddField(field.Name, field.Value, field.Inline)
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embed:      embed.MessageEmbed,
		Components: []discordgo.MessageComponent{actionsRow},
	})

	if err != nil {
		fmt.Printf("Error sending complex notification to channel %s: %v\n", channelId, err)
		return nil, err
	}

	return &sentMessage.ID, nil

}
