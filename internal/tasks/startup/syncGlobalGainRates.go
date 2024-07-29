package startup

import (
	"fmt"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
)

// Updates the runtime application variables with the latest gain rates
// as available in the DB.
func SyncGlobalGainRates() {

	globalGainRates, err := repositories.NewGlobalGainRatesRepository().GetGlobalGainRates()
	if err != nil {
		fmt.Printf("Failed to sync runtime global gain rates with values from DB: %v\n", err)
		return
	}

	for _, rate := range globalGainRates {
		switch rate.ActivityId {
		case "msg_send":
			globalConfiguration.ExperienceReward_MessageSent = rate.MultiplierXp
			globalConfiguration.CoinReward_MessageSent = rate.MultiplierCoins
		case "slash_command_used":
			globalConfiguration.ExperienceReward_SlashCommandUsed = rate.MultiplierXp
			globalConfiguration.CoinReward_SlashCommandUsed = rate.MultiplierCoins
		case "react_recv":
			globalConfiguration.ExperienceReward_ReactionReceived = rate.MultiplierXp
			globalConfiguration.CoinReward_ReactionReceived = rate.MultiplierCoins
		case "in_vc":
			globalConfiguration.ExperienceReward_InVc = rate.MultiplierXp
			globalConfiguration.CoinReward_InVc = rate.MultiplierCoins
		case "in_music":
			globalConfiguration.ExperienceReward_InMusic = rate.MultiplierXp
			globalConfiguration.CoinReward_InMusic = rate.MultiplierCoins
		default:
			fmt.Printf("Activity %s not supported for variable gain rates", rate.ActivityId)
		}
	}

}
