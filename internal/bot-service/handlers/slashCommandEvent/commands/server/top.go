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
	topCount := 12

	embed := embed.NewEmbed().
		SetTitle("🏆   OTA Server Global Leaderboard").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Top by messages sent
	ProcessTopEmbed(topCount, s, i.Interaction, embed)

	globals.LastUsedTopTimestamp = time.Now()

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

			// Process the time spent in VCs in a nice format
			var timeSpentInVcs string = ""
			sTimeSpentInVc := int64(stats.TimeSpentInVoiceChannels)
			daysVC, hoursVC, minutesVC, secondsVC := utils.HumanReadableTimeLength(float64(sTimeSpentInVc))
			if daysVC != 0 {
				timeSpentInVcs = fmt.Sprintf("%dd, %dh:%dm:%ds", daysVC, hoursVC, minutesVC, secondsVC)
			} else if daysVC == 0 && hoursVC != 0 {
				timeSpentInVcs = fmt.Sprintf("%dh:%dm:%ds", hoursVC, minutesVC, secondsVC)
			} else if daysVC == 0 && hoursVC == 0 {
				timeSpentInVcs = fmt.Sprintf("%dm:%ds", minutesVC, secondsVC)
			}

			// Process the time spent listening to music a nice format
			var timeSpentListeningMusic string = ""
			sTimeSpentListeningMusic := int64(stats.TimeSpentListeningToMusic)
			daysMusic, hoursMusic, minutesMusic, secondsMusic := utils.HumanReadableTimeLength(float64(sTimeSpentListeningMusic))
			if daysMusic != 0 {
				timeSpentListeningMusic = fmt.Sprintf("%dd, %dh:%dm:%ds", daysMusic, hoursMusic, minutesMusic, secondsMusic)
			} else if daysMusic == 0 && hoursMusic != 0 {
				timeSpentListeningMusic = fmt.Sprintf("%dh:%dm:%ds", hoursMusic, minutesMusic, secondsMusic)
			} else if daysMusic == 0 && hoursMusic == 0 {
				timeSpentListeningMusic = fmt.Sprintf("%dm:%ds", minutesMusic, secondsMusic)
			}

			rankingRowName := fmt.Sprintf("**%d.** %s**%s**", idx+1, rankMedal, topUser.DiscordTag)
			rankingRowValue := fmt.Sprintf("Total: `%d` XP 💠 | `%d` ✉️ | `%d` 💯 | `%s` 🎙️ | `%s` 🎵\n", int(topUser.XpGained), stats.NumberMessagesSent, stats.NumberReactionsReceived, timeSpentInVcs, timeSpentListeningMusic)
			embed.AddField(rankingRowName, rankingRowValue, false)
		}
		embed.AddField("", topContentText, false)
	}
}
