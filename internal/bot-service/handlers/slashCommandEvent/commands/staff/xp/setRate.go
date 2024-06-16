package xpRateSettingSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/notifications"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
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
	activityName, multiplierName := getArgumentDisplayNames(activity, multiplierStringInput)

	switch activity {
	case "msg_send":
		if multiplierStringInput == "def" {
			globals.ExperienceReward_MessageSent = globals.DefaultExperienceReward_MessageSent
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			globals.ExperienceReward_MessageSent = *multiplier
		}
	case "react_recv":
		if multiplierStringInput == "def" {
			globals.ExperienceReward_ReactionReceived = globals.DefaultExperienceReward_ReactionReceived
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			globals.ExperienceReward_ReactionReceived = *multiplier
		}
	case "slash_use":
		if multiplierStringInput == "def" {
			globals.ExperienceReward_SlashCommandUsed = globals.DefaultExperienceReward_SlashCommandUsed
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			globals.ExperienceReward_SlashCommandUsed = *multiplier
		}
	case "spent_vc":
		if multiplierStringInput == "def" {
			globals.ExperienceReward_InVc = globals.DefaultExperienceReward_InVc
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			globals.ExperienceReward_InVc = *multiplier
		}
	case "spent_music":
		if multiplierStringInput == "def" {
			globals.ExperienceReward_InMusic = globals.DefaultExperienceReward_InMusic
		} else {
			multiplier, convErr := utils.StringToFloat64(multiplierStringInput)
			if convErr != nil {
				errMsg := fmt.Sprintf("The provided `multiplier` command argument is invalid. (term: `%s`)", multiplierName)
				utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
				return
			}
			globals.ExperienceReward_InMusic = *multiplier
		}
	}

	// Send notification to target staff channel to announce the global rate change
	if channel, channelExists := globals.NotificationChannels["notif-aztebot"]; channelExists {
		go sendXpRateChangeNotification(s, channel.ChannelId, activityName, multiplierName)
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   Updated Global XP Rate For `%s`", activityName)).
		SetColor(000000).
		AddField(fmt.Sprintf("New gain rate is `%s`.", multiplierName), "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func sendXpRateChangeNotification(s *discordgo.Session, channelId string, activityName string, multiplierName string) {

	fields := []discordgo.MessageEmbedField{
		{
			Name:   "Activity",
			Value:  activityName,
			Inline: false,
		},
		{
			Name:   "Rate Multiplier",
			Value:  multiplierName,
			Inline: false,
		},
		{
			Name:   "\u200B",
			Value:  "",
			Inline: false,
		},
		{
			Name:   "",
			Value:  "|@everyone|",
			Inline: false,
		},
	}

	notificationTitle := fmt.Sprintf("Global XP Rate Change for `%s`", activityName)
	notifications.SendNotificationToTextChannel(s, channelId, notificationTitle, fields, true)

}

func getArgumentDisplayNames(activityInput string, multiplierInput string) (string, string) {

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
