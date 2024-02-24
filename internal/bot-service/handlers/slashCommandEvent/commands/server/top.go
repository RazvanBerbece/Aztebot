package serverSlashHandlers

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

	durationSinceLastRanksCommand := time.Since(globals.LastUsedTopTimestamp)
	if int(durationSinceLastRanksCommand.Minutes()) < 5 {
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
	topCount := 15

	embed := embed.NewEmbed().
		SetTitle("🏆   OTA Server Global Leaderboard").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Top by messages sent
	ProcessTopEmbed(topCount, s, i.Interaction, embed)

	globals.LastUsedRanksTimestamp = time.Now()

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}

func ProcessTopEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {

	topXpGains, err := globalsRepo.UserStatsRepository.GetTopUsersByXp(topCount)
	if err != nil {
		log.Printf("Cannot retrieve global OTA leaderboard: %v", err)
	}

	if len(topXpGains) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topXpGains {
			// Dynamically add a medal emoji depending on position in ranking
			rankMedal := ""
			switch idx {
			case 0:
				rankMedal = "🥇 "
			case 1:
				rankMedal = "🥈 "
			case 2:
				rankMedal = "🥉 "
			default:
				rankMedal = ""
			}

			// Get rest of stats for user to display in the result embed
			stats, err := globalsRepo.UserStatsRepository.GetStatsForUser(topUser.UserId)
			if err != nil {
				log.Printf("Cannot retrieve stats for user: %v", err)
				continue
			}

			rankingRowName := fmt.Sprintf("**%d.** %s**%s**", idx+1, rankMedal, topUser.DiscordTag)
			rankingRowValue := fmt.Sprintf("Total: `%d` XP 💠 | `%d` ✉️ | `%d` 💯 | `%d` 🎙️ | `%d` 🎵\n", int(topUser.XpGained), stats.NumberMessagesSent, stats.NumberReactionsReceived, stats.TimeSpentInVoiceChannels, stats.TimeSpentListeningToMusic)
			embed.AddField(rankingRowName, rankingRowValue, false)
		}
		embed.AddField("", topContentText, false)
	}
}
