package utils

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func ErrorEmbedResponseEdit(s *discordgo.Session, i *discordgo.Interaction, errorMessage string) {

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖❌   `/%s` Command Execution Error", i.ApplicationCommandData().Name)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Error Report", errorMessage, false)

	editWebhook := discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}

	s.InteractionResponseEdit(i, &editWebhook)
}

func SendErrorEmbedResponse(s *discordgo.Session, i *discordgo.Interaction, errorMessage string) {

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖❌   `/%s` Command Execution Error", i.ApplicationCommandData().Name)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Error Report", errorMessage, false)

	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
