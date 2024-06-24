package member

import (
	"database/sql"
	"fmt"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	aztemarketRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztemarket"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/bwmarrin/discordgo"
)

// By default logs errors and state to Discord.
func AwardFunds(s *discordgo.Session, usersRepository repositories.UsersRepository, walletsRepository aztemarketRepositories.WalletsRepository, userId string, funds float64, activity string) error {

	if funds < 0 || funds > 500000.0 {
		return fmt.Errorf("cannot award funds to user with ID `%s`, because the number of awarded `funds` (`%.2f`) is invalid", userId, funds)
	}

	err := walletsRepository.AddFundsToWalletForUser(userId, funds)
	if err != nil {
		log := fmt.Sprintf("An error ocurred while awarding funds to user `%s`: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
		go discordChannelLogger.LogError(log)
		return err
	}

	walletId, err := walletsRepository.GetWalletIdForUser(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			// user doesn't have a wallet currently. don't log to ensure that logs stay relatively noise free
			return nil
		}
		log := fmt.Sprintf("An error ocurred while retrieving wallet ID for user `%s`: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogError(log)
		return err
	}

	user, err := usersRepository.GetUser(userId)
	if err != nil {
		log := fmt.Sprintf("An error ocurred while retrieving user `%s` to awards funds to: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogError(log)
		return err
	}

	// Add audit log to ledger channel to keep a track record of *all* coin awards
	log := fmt.Sprintf("Awarded `🪙 %.2f` AzteCoins\nto user `%s` (`%s`)\nwith wallet ID `%s`\nfor activity ID `%s`", funds, user.DiscordTag, userId, *walletId, activity)
	discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
	go discordChannelLogger.LogInfo(log)

	return nil

}
