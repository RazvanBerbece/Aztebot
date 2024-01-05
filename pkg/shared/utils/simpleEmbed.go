package utils

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func SimpleEmbed(title string, description string) []*discordgo.MessageEmbed {
	embed := embed.NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
