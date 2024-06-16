package channelHandlers

import (
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandleCoinAwardEvents(s *discordgo.Session, logger logging.Logger) {

	for coinAwardEvent := range globalMessaging.CoinAwardsChannel {
		go member.AwardFunds(s, coinAwardEvent.UserId, coinAwardEvent.Funds, coinAwardEvent.Activity)
	}

}
