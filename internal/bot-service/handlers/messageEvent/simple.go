package messageEvent

import (
	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/logger"
	"github.com/bwmarrin/discordgo"
)

func SimpleMsgReply(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "robotel" {
		logger.LogHandlerCall("SimpleMsgReply (robotel)", "")
		s.ChannelMessageSend(m.ChannelID, "Prezent! Cu ce te pot ajuta?")
	}
	if m.Content == "mergi?" {
		logger.LogHandlerCall("SimpleMsgReply (mergi?)", "")
		s.ChannelMessageSend(m.ChannelID, "Dupa cum se vede, sunt activ si raspund la comenzi!")
	}
	if m.Content == "cat e ceasul?" {
		logger.LogHandlerCall("SimpleMsgReply (cat e ceasul?)", "")
		s.ChannelMessageSend(m.ChannelID, "Cat ti-e nasul, hahaha!")
	}

}
