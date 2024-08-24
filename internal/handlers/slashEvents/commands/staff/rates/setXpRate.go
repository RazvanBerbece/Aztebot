package gainRatesSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashSetGlobalXpRateForActivity(s *discordgo.Session, i *discordgo.InteractionCreate) {

	activity := i.ApplicationCommandData().Options[0].StringValue()
	multiplierStringInput := i.ApplicationCommandData().Options[1].StringValue()

	// Dirty Hack 25 Feb 2024
	// It seems that it's not straightforward at all to get the display name of the argument option,
	// so we resort to this for the meantime to get a nicely looking activity and multiplier name
	activityName, multiplierName := GetArgumentDisplayNames(activity, multiplierStringInput)

	switch activity {
	case "msg_send":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_MessageSend,
				globalConfiguration.DefaultExperienceReward_MessageSent)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_MessageSend, err)
				return
			}
			globalConfiguration.ExperienceReward_MessageSent = globalConfiguration.DefaultExperienceReward_MessageSent
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_MessageSend,
				globalConfiguration.DefaultExperienceReward_MessageSent**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_MessageSend, err)
				return
			}
			globalConfiguration.ExperienceReward_MessageSent = globalConfiguration.DefaultExperienceReward_MessageSent * *multiplier
		}
	case "react_recv":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_ReactionReceived,
				globalConfiguration.DefaultExperienceReward_ReactionReceived)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_ReactionReceived, err)
				return
			}
			globalConfiguration.ExperienceReward_ReactionReceived = globalConfiguration.DefaultExperienceReward_ReactionReceived
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_ReactionReceived,
				globalConfiguration.DefaultExperienceReward_ReactionReceived**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_ReactionReceived, err)
				return
			}
			globalConfiguration.ExperienceReward_ReactionReceived = globalConfiguration.DefaultExperienceReward_ReactionReceived * *multiplier
		}
	case "slash_use":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_SlashCommandUse,
				globalConfiguration.DefaultExperienceReward_SlashCommandUsed)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_SlashCommandUse, err)
				return
			}
			globalConfiguration.ExperienceReward_SlashCommandUsed = globalConfiguration.DefaultExperienceReward_SlashCommandUsed
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_SlashCommandUse,
				globalConfiguration.DefaultExperienceReward_SlashCommandUsed**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_SlashCommandUse, err)
				return
			}
			globalConfiguration.ExperienceReward_SlashCommandUsed = globalConfiguration.DefaultExperienceReward_SlashCommandUsed * *multiplier
		}
	case "spent_vc":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_TimeInVc,
				globalConfiguration.DefaultExperienceReward_InVc)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInVc, err)
				return
			}
			globalConfiguration.ExperienceReward_InVc = globalConfiguration.DefaultExperienceReward_InVc
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_TimeInVc,
				globalConfiguration.DefaultExperienceReward_InVc**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInVc, err)
				return
			}
			globalConfiguration.ExperienceReward_InVc = globalConfiguration.DefaultExperienceReward_InVc * *multiplier
		}
	case "spent_music":
		if multiplierStringInput == "def" {
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_TimeInMusic,
				globalConfiguration.DefaultExperienceReward_InMusic)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInMusic, err)
				return
			}
			globalConfiguration.ExperienceReward_InMusic = globalConfiguration.DefaultExperienceReward_InMusic
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			err := repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
				globalConfiguration.ActivityId_TimeInMusic,
				globalConfiguration.DefaultExperienceReward_InMusic**multiplier)
			if err != nil {
				fmt.Printf("Failed to update global XP gain rate for %s: %v\n", globalConfiguration.ActivityId_TimeInMusic, err)
				return
			}
			globalConfiguration.ExperienceReward_InMusic = globalConfiguration.DefaultExperienceReward_InMusic * *multiplier
		}
	}

	// Send notification to target staff channel to announce the global rate change
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-aztebotUpdatesChannel"]; channelExists {
		go SendGlobalRateChangeNotification(channel.ChannelId, "💠 XP", activityName, multiplierName)
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖   Updated Global XP Rate For `%s`", activityName)).
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

func SendGlobalRateChangeNotification(channelId string, rateName string, activityName string, multiplierName string) {

	// Build global XP rate change embed
	embed := embed.
		NewEmbed().
		SetAuthor("AzteBot Global Broadcaster", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖📣	Gain Rate Boost Announcement - `%s`", rateName)).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000)

	if multiplierName == "Default OTA Value" {
		embed.AddField("", fmt.Sprintf("Activities involving `%s` are now worth the default amount of `%s`.", activityName, rateName), false)
	} else {
		embed.AddField("", fmt.Sprintf("Activities involving `%s` are now worth `%s` as many `%s` !", activityName, multiplierName, rateName), false)
	}

	embed.AtTagEveryone(true)

	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           embed,
	}

}

func GetArgumentDisplayNames(activityInput string, multiplierInput string) (string, string) {

	var activityName string
	var multiplierName string

	switch activityInput {
	case "msg_send":
		activityName = "Message Sends"
	case "react_recv":
		activityName = "Reactions Received"
	case "slash_use":
		activityName = "Slash Commands Used"
	case "spent_vc":
		activityName = "Time Spent in Voice Channels"
	case "spent_music":
		activityName = "Time Spent Listening to Music"
	}

	switch multiplierInput {
	case "def":
		multiplierName = "Default OTA Value"
	case "1.5":
		multiplierName = "1.5x"
	case "2.0":
		multiplierName = "2x"
	case "3.0":
		multiplierName = "3x"
	}

	return activityName, multiplierName
}
