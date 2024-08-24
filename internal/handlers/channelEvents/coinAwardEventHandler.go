package channelHandlers

import (
	aztebotRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	aztemarketRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztemarket"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

// Use separate repositories (and connections) for the coin awards,
// as there can be many awards due to voice tick rewards, hence lots of egress traffic to the MySQL instance

var coinAwardSpecific_UsersRepository = aztebotRepositories.NewUsersRepository()
var coinAwardSpecific_walletsRepository = aztemarketRepositories.NewWalletsRepository(globalConfiguration.MySqlAztemarketRootConnectionString)

func HandleCoinAwardEvents(s *discordgo.Session, logger logging.Logger) {

	for coinAwardEvent := range globalMessaging.CoinAwardsChannel {
		go member.AwardFunds(s, coinAwardEvent.GuildId, *coinAwardSpecific_UsersRepository, coinAwardSpecific_walletsRepository, coinAwardEvent.UserId, coinAwardEvent.Funds, coinAwardEvent.Activity)
	}

}
