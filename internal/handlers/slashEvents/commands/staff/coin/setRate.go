package coinSlashHandlers

import (
	"fmt"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	xpSystemSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/xp"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashSetGlobalCoinRateForActivity(s *discordgo.Session, i *discordgo.InteractionCreate) {

	activity := i.ApplicationCommandData().Options[0].StringValue()
	multiplierStringInput := i.ApplicationCommandData().Options[1].StringValue()

	// Dirty Hack 14 Jun 2024 (still the case...)
	activityName, multiplierName := xpSystemSlashHandlers.GetArgumentDisplayNames(activity, multiplierStringInput)

	switch activity {
	case "msg_send":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_MessageSend,
				globalConfiguration.DefaultCoinReward_MessageSent)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_MessageSend, err)
				return
			}
			globalConfiguration.CoinReward_MessageSent = globalConfiguration.DefaultCoinReward_MessageSent
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_MessageSend,
				globalConfiguration.DefaultCoinReward_MessageSent**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_MessageSend, err)
				return
			}
			globalConfiguration.CoinReward_MessageSent = globalConfiguration.DefaultCoinReward_MessageSent * *multiplier
		}
	case "react_recv":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_ReactionReceived,
				globalConfiguration.DefaultCoinReward_ReactionReceived)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_ReactionReceived, err)
				return
			}
			globalConfiguration.CoinReward_ReactionReceived = globalConfiguration.DefaultCoinReward_ReactionReceived
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_ReactionReceived,
				globalConfiguration.DefaultCoinReward_ReactionReceived**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_ReactionReceived, err)
				return
			}
			globalConfiguration.CoinReward_ReactionReceived = globalConfiguration.DefaultCoinReward_ReactionReceived * *multiplier
		}
	case "slash_use":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_SlashCommandUse,
				globalConfiguration.DefaultCoinReward_SlashCommandUsed)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_SlashCommandUse, err)
				return
			}
			globalConfiguration.CoinReward_SlashCommandUsed = globalConfiguration.DefaultCoinReward_SlashCommandUsed
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_SlashCommandUse,
				globalConfiguration.DefaultCoinReward_SlashCommandUsed**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_SlashCommandUse, err)
				return
			}
			globalConfiguration.CoinReward_SlashCommandUsed = globalConfiguration.DefaultCoinReward_SlashCommandUsed * *multiplier
		}
	case "spent_vc":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_TimeInVc,
				globalConfiguration.DefaultCoinReward_InVc)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInVc, err)
				return
			}
			globalConfiguration.CoinReward_InVc = globalConfiguration.DefaultCoinReward_InVc
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_TimeInVc,
				globalConfiguration.DefaultCoinReward_InVc**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInVc, err)
				return
			}
			globalConfiguration.CoinReward_InVc = globalConfiguration.DefaultCoinReward_InVc * *multiplier
		}
	case "spent_music":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_TimeInMusic,
				globalConfiguration.DefaultCoinReward_InMusic)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInMusic, err)
				return
			}
			globalConfiguration.CoinReward_InMusic = globalConfiguration.DefaultCoinReward_InMusic
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
				globalConfiguration.ActivityId_TimeInMusic,
				globalConfiguration.DefaultCoinReward_InMusic**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global coins gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInMusic, err)
				return
			}
			globalConfiguration.CoinReward_InMusic = globalConfiguration.DefaultCoinReward_InMusic * *multiplier
		}
	}

	// Send notification to target staff channel to announce the global rate change
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-aztebotUpdatesChannel"]; channelExists {
		go xpSystemSlashHandlers.SendGlobalRateChangeNotification(channel.ChannelId, "🪙 AzteCoins", activityName, multiplierName)
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖   Updated Global Coin Gain Rate For `%s`", activityName)).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000).
		AddField(fmt.Sprintf("New gain rate is `%s`.", multiplierName), "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}
