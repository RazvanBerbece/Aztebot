package utils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func DeleteInteractionResponse(s *discordgo.Session, i *discordgo.Interaction, msDelay int) {

	time.Sleep(time.Duration(msDelay) * time.Millisecond)

	// Delete the response
	s.InteractionResponseDelete(i)
}
