package utils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func DeleteInteractionResponse(s *discordgo.Session, i *discordgo.Interaction, delay int) {

	time.Sleep(time.Duration(delay) * time.Second)

	// Delete the response
	s.InteractionResponseDelete(i)
}
