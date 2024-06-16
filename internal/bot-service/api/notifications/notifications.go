package notifications

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func SendNotificationToTextChannel(s *discordgo.Session, channelId string, notificationTitle string, fields []discordgo.MessageEmbedField) error {

	// Build notification embed
	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖ℹ️   %s", notificationTitle)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)
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
