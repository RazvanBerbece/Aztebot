package messageEvent

import (
	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/logger"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		logger.LogHandlerCall("ping", "")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

}
