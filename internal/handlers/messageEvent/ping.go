package messageEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		utils.LogHandlerCall("ping", "")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

}
