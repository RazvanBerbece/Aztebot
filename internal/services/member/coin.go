package member

import (
	"database/sql"
	"fmt"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	aztemarketRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztemarket"
	"github.com/RazvanBerbece/Aztebot/internal/services/economy"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/bwmarrin/discordgo"
)

// By default logs errors and state to Discord.
func AwardFunds(s *discordgo.Session, guildId string, economicsRepository aztemarketRepositories.DbCurrencySystemStateRepository, usersRepository repositories.UsersRepository, walletsRepository aztemarketRepositories.WalletsRepository, userId string, funds float64, activity string) error {

	_, err := walletsRepository.GetWalletIdForUser(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			// user doesn't have a wallet currently. return early
			return nil
		}
		log := fmt.Sprintf("An error ocurred while retrieving awardee wallet ID for user `%s`: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogError(log)
		return err
	}

	if funds < 0 || funds > 500000.0 {
		return fmt.Errorf("cannot award funds to user with ID `%s`, because the number of awarded `funds` (`%.2f`) is invalid", userId, funds)
	}

	// Check that the economy pot permits the coin award
	economicSystem, err := economicsRepository.GetCurrencyStateForGuild(guildId)
	if err != nil {
		if err == sql.ErrNoRows {
			log := "No economic system available for this guild. Make sure you create one before awarding currency for activities by using the `/economy-create` command."
			discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
			go discordChannelLogger.LogError(log)
			return err
		}
		log := fmt.Sprintf("An error ocurred while retrieving economic system for fund awarding to user `%s`: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
		go discordChannelLogger.LogError(log)
		return err
	}

	if economicSystem.TotalCurrencyAvailable < funds {
		log := fmt.Sprintf("Currency `%s` ran out of available global units for activity awards", economicSystem.CurrencyName)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
		go discordChannelLogger.LogError(log)
		return err
	}

	err = walletsRepository.AddFundsToWalletForUser(userId, funds)
	if err != nil {
		log := fmt.Sprintf("An error ocurred while awarding funds to user `%s`: %v\n", userId, err)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
		go discordChannelLogger.LogError(log)
		return err
	}

	economyService := economy.EconomyService{
		CurrencySystemStateRepository: economicsRepository,
	}

	err = economyService.AllocateFlowingCurrencyForGuild(guildId, funds)
	if err != nil {
		// Rollback
		rollbackErr := walletsRepository.SubtractFundsFromWallet(userId, funds)
		if rollbackErr != nil {
			log := fmt.Sprintf("An error ocurred while subtracting funds from user `%s`: %v\n", userId, rollbackErr)
			discordChannelLogger := logging.NewDiscordLogger(s, "notif-economyDebug")
			go discordChannelLogger.LogError(log)
			return err
		}
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-economyDebug")
		go discordChannelLogger.LogError(err.Error())
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
	log := fmt.Sprintf("Awarded `🪙 %.2f` AzteCoins\nto user `%s` (`%s`)\nfor activity ID `%s`", funds, user.DiscordTag, userId, activity)
	discordChannelLogger := logging.NewDiscordLogger(s, "notif-coinAwards")
	go discordChannelLogger.LogInfo(log)

	return nil

}
