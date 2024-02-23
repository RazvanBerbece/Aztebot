package messageEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SimpleMsgReply(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "robotel" {
		utils.LogHandlerCall("SimpleMsgReply (robotel)", "")
		s.ChannelMessageSend(m.ChannelID, "Prezent! Cu ce te pot ajuta?")
	}
	if m.Content == "mergi?" {
		utils.LogHandlerCall("SimpleMsgReply (mergi?)", "")
		s.ChannelMessageSend(m.ChannelID, "Dupa cum se vede, sunt activ si raspund la comenzi!")
	}
	if m.Content == "cat e ceasul?" {
		utils.LogHandlerCall("SimpleMsgReply (cat e ceasul?)", "")
		s.ChannelMessageSend(m.ChannelID, "Cat ti-e nasul, hahaha!")
	}

}
