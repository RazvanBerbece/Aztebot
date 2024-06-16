package slashHandlers

import (
	"fmt"
	"log"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTop(s *discordgo.Session, i *discordgo.InteractionCreate) {

	go processTopCommand(s, i)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Processing `/top` command ...",
		},
	})
}

func processTopCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

	embed := embed.NewEmbed().
		SetTitle("🤖   OTA Server Leaderboard").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Top by messages sent
	topCount := 5
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
		for idx, topUser := range topMessagesSent {
			embed.AddField("", fmt.Sprintf("**%d.** **%s**    (sent `%d` ✉️)", idx+1, topUser.DiscordTag, topUser.MessagesSent), false)
		}
	}

	// Top by time spent in VCs
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
		for idx, topUser := range topTimeInVCs {
			days, hours, minutes, seconds := utils.HumanReadableTimeLength(float64(topUser.TimeSpentInVCs))
			embed.AddField("", fmt.Sprintf("**%d.** **%s** (spent `%dd, %dh:%dm:%ds` in voice channels 🎙️)", idx+1, topUser.DiscordTag, days, hours, minutes, seconds), false)
		}
	}

	embeds := []*discordgo.MessageEmbed{embed.MessageEmbed}

	// The edit webhook container holds the updated interaction response details (contents, embeds, etc.)
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &embeds,
	}

	s.InteractionResponseEdit(i.Interaction, &editWebhook)
}
