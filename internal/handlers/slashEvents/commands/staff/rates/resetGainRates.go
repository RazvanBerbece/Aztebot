package gainRatesSlashHandlers

import (
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashResetGainRates(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// COIN AWARDS
	go repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
		globalConfiguration.ActivityId_MessageSend,
		globalConfiguration.DefaultCoinReward_MessageSent)

	go repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
		globalConfiguration.ActivityId_ReactionReceived,
		globalConfiguration.DefaultCoinReward_ReactionReceived)

	go repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
		globalConfiguration.ActivityId_SlashCommandUse,
		globalConfiguration.DefaultCoinReward_SlashCommandUsed)

	go repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
		globalConfiguration.ActivityId_TimeInVc,
		globalConfiguration.DefaultCoinReward_InVc)

	go repositories.NewGlobalGainRatesRepository().UpdateCoinsGlobalGainRate(
		globalConfiguration.ActivityId_TimeInMusic,
		globalConfiguration.DefaultCoinReward_InMusic)

	// XP AWARDS
	go repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
		globalConfiguration.ActivityId_MessageSend,
		globalConfiguration.DefaultExperienceReward_MessageSent)

	go repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
		globalConfiguration.ActivityId_ReactionReceived,
		globalConfiguration.DefaultExperienceReward_ReactionReceived)

	go repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
		globalConfiguration.ActivityId_SlashCommandUse,
		globalConfiguration.DefaultExperienceReward_SlashCommandUsed)

	go repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
		globalConfiguration.ActivityId_TimeInVc,
		globalConfiguration.DefaultExperienceReward_InVc)

	go repositories.NewGlobalGainRatesRepository().UpdateXpGlobalGainRate(
		globalConfiguration.ActivityId_TimeInMusic,
		globalConfiguration.DefaultExperienceReward_InMusic)

	// Send notification to target staff channel to announce the global rate change
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-aztebotUpdatesChannel"]; channelExists {
		go SendGlobalRateChangeResetNotification(channel.ChannelId)
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("🤖   Reset Global Coin Gain Rates").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000).
		AddField("New gain rates are back to their default values.", "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func SendGlobalRateChangeResetNotification(channelId string) {

	// Build global XP rate change embed
	embed := embed.
		NewEmbed().
		SetAuthor("AzteBot Global Broadcaster", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("🤖📣	Gain Rate Reset Announcement").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000)

	embed.AddField("", "All activities will now award the default amount of `🪙 AzteCoin` and `💠 XP`.", false)

	embed.AtTagEveryone(true)

	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           embed,
	}

}
