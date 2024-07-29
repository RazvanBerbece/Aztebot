package serverSlashHandlers

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashServerGainRates(s *discordgo.Session, i *discordgo.InteractionCreate) {

	embedToSend := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("🤖   Server Gain Rates").
		SetColor(000000)

	// Add fields per gain rate
	embedToSend.AddField(
		"Message Sends",
		fmt.Sprintf(
			"`💠 XP` Gain Rate: `%.2f`\n`🪙 AzteCoins` Gain Rate: `%.2f`",
			globalConfiguration.ExperienceReward_MessageSent, globalConfiguration.CoinReward_MessageSent),
		false)

	embedToSend.AddField(
		"Slash Command Usage",
		fmt.Sprintf(
			"`💠 XP` Gain Rate: `%.2f`\n`🪙 AzteCoins` Gain Rate: `%.2f`",
			globalConfiguration.ExperienceReward_SlashCommandUsed, globalConfiguration.CoinReward_SlashCommandUsed),
		false)

	embedToSend.AddField(
		"Reaction Received",
		fmt.Sprintf(
			"`💠 XP` Gain Rate: `%.2f`\n`🪙 AzteCoins` Gain Rate: `%.2f`",
			globalConfiguration.ExperienceReward_ReactionReceived, globalConfiguration.CoinReward_ReactionReceived),
		false)

	embedToSend.AddField(
		"Time Spent in Voice Channels",
		fmt.Sprintf(
			"`💠 XP` Gain Rate: `%.5f`\n`🪙 AzteCoins` Gain Rate: `%.5f`",
			globalConfiguration.ExperienceReward_InVc, globalConfiguration.CoinReward_InVc),
		false)

	embedToSend.AddField(
		"Time Spent in Music Channels",
		fmt.Sprintf(
			"`💠 XP` Gain Rate: `%.5f`\n`🪙 AzteCoins` Gain Rate: `%.5f`",
			globalConfiguration.ExperienceReward_InMusic, globalConfiguration.CoinReward_InMusic),
		false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embedToSend.MessageEmbed},
		},
	})
}
