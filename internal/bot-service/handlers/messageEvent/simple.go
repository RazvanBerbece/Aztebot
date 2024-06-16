package messageEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func SimpleMsgReply(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "robotel" {
		logging.LogHandlerCall("SimpleMsgReply (robotel)", "")
		s.ChannelMessageSend(m.ChannelID, "Prezent! Cu ce te pot ajuta?")
	}
	if m.Content == "mergi?" {
		logging.LogHandlerCall("SimpleMsgReply (mergi?)", "")
		s.ChannelMessageSend(m.ChannelID, "Dupa cum se vede, sunt activ si raspund la comenzi!")
	}
	if m.Content == "cat e ceasul?" {
		logging.LogHandlerCall("SimpleMsgReply (cat e ceasul?)", "")
		s.ChannelMessageSend(m.ChannelID, "Cat ti-e nasul, hahaha!")
	}

}
