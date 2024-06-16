package messageEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		logging.LogHandlerCall("ping", "")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

}
