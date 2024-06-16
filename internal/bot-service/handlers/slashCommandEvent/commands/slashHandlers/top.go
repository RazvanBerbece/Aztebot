package slashHandlers

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTop(s *discordgo.Session, i *discordgo.InteractionCreate) {

	durationSinceLastTopCommand := time.Since(globals.LastUsedTopTimestamp)
	if int(durationSinceLastTopCommand.Minutes()) < 5 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SimpleEmbed("🤖   Slash Command Usage Limit", "The `/top` slash command can be used only once every 5 minutes to reduce the resource usage of the `AzteBot`."),
			},
		})
		return
	}

	// Initial response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/top` command..."),
		},
	})

	// Final response
	results := TopCommandResultsEmbed(s, i)
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &results,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}

func TopCommandResultsEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) []*discordgo.MessageEmbed {

	// Leaderboard parameterisation
	topCount := 5

	embed := embed.NewEmbed().
		SetTitle("🤖   OTA Server Leaderboard").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Top by messages sent
	ProcessTopMessagesPartialEmbed(topCount, s, i.Interaction, embed)
	// updateInteraction(s, i.Interaction, embed)

	// Top by time spent in VCs
	ProcessTopVCSpentPartialEmbed(topCount, s, i.Interaction, embed)
	// updateInteraction(s, i.Interaction, embed)

	// Top by active day streak
	ProcessTopActiveDayStreakPartialEmbed(topCount, s, i.Interaction, embed)
	// updateInteraction(s, i.Interaction, embed)

	// Top by reactions received
	ProcessTopReactionsReceivedPartialEmbed(topCount, s, i.Interaction, embed)
	// updateInteraction(s, i.Interaction, embed)

	globals.LastUsedTopTimestamp = time.Now()

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}

func ProcessTopMessagesPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topMessagesSent, err := globalsRepo.UserStatsRepository.GetTopUsersByMessageSent(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top messages sent from the Discord API: %v", err)
	}
	embed.
		AddLineBreakField().
		AddField(fmt.Sprintf("✉️ Top %d By Messages Sent", topCount), "", false)
	if len(topMessagesSent) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topMessagesSent {
			topContentText += fmt.Sprintf("**%d.** **%s**    (sent `%d` ✉️)\n", idx+1, topUser.DiscordTag, topUser.MessagesSent)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopVCSpentPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topTimeInVCs, err := globalsRepo.UserStatsRepository.GetTopUsersByTimeSpentInVC(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top times spent in VC from the Discord API: %v", err)
	}
	embed.
		AddLineBreakField().
		AddField(fmt.Sprintf("🎙️ Top %d By Time Spent in Voice Channels", topCount), "", false)
	if len(topTimeInVCs) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topTimeInVCs {
			days, hours, minutes, seconds := utils.HumanReadableTimeLength(float64(topUser.TimeSpentInVCs))
			topContentText += fmt.Sprintf("**%d.** **%s** (spent `%dd, %dh:%dm:%ds` in voice channels 🎙️)\n", idx+1, topUser.DiscordTag, days, hours, minutes, seconds)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopActiveDayStreakPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topStreaks, err := globalsRepo.UserStatsRepository.GetTopUsersByActiveDayStreak(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top streaks from the Discord API: %v", err)
	}
	embed.
		AddLineBreakField().
		AddField(fmt.Sprintf("🔄 Top %d By Active Day Streak", topCount), "", false)
	if len(topStreaks) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topStreaks {
			topContentText += fmt.Sprintf("**%d.** **%s** (active for `%d` days in a row 🔄)\n", idx+1, topUser.DiscordTag, topUser.Streak)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopReactionsReceivedPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topReactions, err := globalsRepo.UserStatsRepository.GetTopUsersByReceivedReactions(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top reactions received from the Discord API: %v", err)
	}
	embed.
		AddLineBreakField().
		AddField(fmt.Sprintf("💯 Top %d By Total Reactions Received", topCount), "", false)
	if len(topReactions) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topReactions {
			topContentText += fmt.Sprintf("**%d.** **%s** (received a total of `%d` reactions 💯)\n", idx+1, topUser.DiscordTag, topUser.ReactionsReceived)
		}
		embed.AddField("", topContentText, false)
	}
}
