package slashHandlers

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func HandleSlashMusicPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s Pong!", os.Getenv("APP_NAME")),
		},
	})
}
