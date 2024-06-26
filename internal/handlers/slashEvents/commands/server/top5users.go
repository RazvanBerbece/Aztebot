package serverSlashHandlers

import (
	"fmt"
	"log"
	"time"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTop5Users(s *discordgo.Session, i *discordgo.InteractionCreate) {

	durationSinceLastTop5sCommand := time.Since(globalState.LastUsedTop5sTimestamp)
	if int(durationSinceLastTop5sCommand.Minutes()) < 5 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SimpleEmbed("🤖   Slash Command Usage Limit", "The `/top5user` slash command can be used only once every 5 minutes to reduce the resource usage of the `AzteBot`."),
			},
		})
		return
	}

	// Initial response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/top5user` command..."),
		},
	})

	// Final response
	results := Top5CommandResultsEmbed(s, i)
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &results,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}

func Top5CommandResultsEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) []*discordgo.MessageEmbed {

	// Leaderboard parameterisation
	topCount := 5

	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("🤖   OTA Server Top 5s Leaderboard").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000)

	// Top by messages sent
	ProcessTopMessagesPartialEmbed(topCount, s, i.Interaction, embed)

	// Top by time spent in VCs
	ProcessTopVCSpentPartialEmbed(topCount, s, i.Interaction, embed)

	// Top by active day streak
	ProcessTopActiveDayStreakPartialEmbed(topCount, s, i.Interaction, embed)

	// Top by reactions received
	ProcessTopReactionsReceivedPartialEmbed(topCount, s, i.Interaction, embed)

	// Top by time spent listening to music
	ProcessTopMusicListeningTimePartialEmbed(topCount, s, i.Interaction, embed)

	globalState.LastUsedTop5sTimestamp = time.Now()

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}

func ProcessTopMessagesPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topMessagesSent, err := globalRepositories.UserStatsRepository.GetTopUsersByMessageSent(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top messages sent from the Discord API: %v", err)
	}
	embed.
		AddField(fmt.Sprintf("✉️ Top %d By Messages Sent", topCount), "", false)
	if len(topMessagesSent) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topMessagesSent {
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
			topContentText += fmt.Sprintf("**%d.** %s**%s**   (sent `%d` ✉️)\n", idx+1, rankMedal, topUser.DiscordTag, topUser.MessagesSent)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopVCSpentPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topTimeInVCs, err := globalRepositories.UserStatsRepository.GetTopUsersByTimeSpentInVC(topCount)
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
			days, hours, minutes, seconds := utils.HumanReadableDuration(float64(topUser.TimeSpentInVCs))
			topContentText += fmt.Sprintf("**%d.** %s**%s** (spent `%dd, %dh:%dm:%ds` in voice channels 🎙️)\n", idx+1, rankMedal, topUser.DiscordTag, days, hours, minutes, seconds)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopActiveDayStreakPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topStreaks, err := globalRepositories.UserStatsRepository.GetTopUsersByActiveDayStreak(topCount)
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
			topContentText += fmt.Sprintf("**%d.** %s**%s** (active for `%d` days in a row 🔄)\n", idx+1, rankMedal, topUser.DiscordTag, topUser.Streak)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopReactionsReceivedPartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topReactions, err := globalRepositories.UserStatsRepository.GetTopUsersByReceivedReactions(topCount)
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
			topContentText += fmt.Sprintf("**%d.** %s**%s** (received a total of `%d` reactions 💯)\n", idx+1, rankMedal, topUser.DiscordTag, topUser.ReactionsReceived)
		}
		embed.AddField("", topContentText, false)
	}
}

func ProcessTopMusicListeningTimePartialEmbed(topCount int, s *discordgo.Session, i *discordgo.Interaction, embed *embed.Embed) {
	topMusicListeners, err := globalRepositories.UserStatsRepository.GetTopUsersByTimeSpentListeningMusic(topCount)
	if err != nil {
		log.Printf("Cannot retrieve OTA leaderboard top times spent listening music from the Discord API: %v", err)
	}
	embed.
		AddLineBreakField().
		AddField(fmt.Sprintf("🎵 Top %d By Time Spent Listening To Music", topCount), "", false)
	if len(topMusicListeners) == 0 {
		embed.AddField("", "No members in this category", false)
	} else {
		topContentText := ""
		for idx, topUser := range topMusicListeners {
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
			days, hours, minutes, seconds := utils.HumanReadableDuration(float64(topUser.TimeSpentListeningMusic))
			topContentText += fmt.Sprintf("**%d.** %s**%s** (spent `%dd, %dh:%dm:%ds` listening to music 🎵)\n", idx+1, rankMedal, topUser.DiscordTag, days, hours, minutes, seconds)
		}
		embed.AddField("", topContentText, false)
	}
}
