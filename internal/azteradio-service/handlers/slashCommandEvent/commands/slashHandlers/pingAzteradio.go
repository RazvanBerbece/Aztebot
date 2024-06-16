package slashHandlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleSlashRadioPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Azteradio Pong!",
		},
	})
}
