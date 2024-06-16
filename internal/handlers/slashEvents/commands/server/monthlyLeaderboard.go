package serverSlashHandlers

import (
	"fmt"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models/dax"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMonthlyLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {

	durationSinceLastMlbCommand := time.Since(globalState.LastUsedMonthlyLeaderboardTimestamp)
	if int(durationSinceLastMlbCommand.Minutes()) < 5 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.SimpleEmbed("🤖   Slash Command Usage Limit", "The `/monthly-leaderboard` slash command can be used only once every 5 minutes to reduce the resource usage of the `AzteBot`."),
			},
		})
		return
	}

	// Initial response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/monthly-leaderboard` command..."),
		},
	})

	// Final response
	results := MonthlyLeaderboardCommandResultsEmbed(s, i)
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &results,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

	globalState.LastUsedMonthlyLeaderboardTimestamp = time.Now()

}

func MonthlyLeaderboardCommandResultsEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) []*discordgo.MessageEmbed {

	maleEntries, err := globalRepositories.MonthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(0)
	if err != nil {
		fmt.Printf("An error ocurred in MonthlyLeaderboardCommandResultsEmbed: %v\n", err)
	}
	femaleEntries, err := globalRepositories.MonthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(1)
	if err != nil {
		fmt.Printf("An error ocurred in MonthlyLeaderboardCommandResultsEmbed: %v\n", err)
	}
	nonbinaryEntries, err := globalRepositories.MonthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(2)
	if err != nil {
		fmt.Printf("An error ocurred in MonthlyLeaderboardCommandResultsEmbed: %v\n", err)
	}
	otherEntries, err := globalRepositories.MonthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(3)
	if err != nil {
		fmt.Printf("An error ocurred in MonthlyLeaderboardCommandResultsEmbed: %v\n", err)
	}

	var kingEntry *dataModels.MonthlyLeaderboardEntry = nil
	var queenEntry *dataModels.MonthlyLeaderboardEntry = nil
	var nonbinaryEntry *dataModels.MonthlyLeaderboardEntry = nil
	var otherEntry *dataModels.MonthlyLeaderboardEntry = nil
	if len(maleEntries) > 0 {
		kingEntry = &maleEntries[0]
	}
	if len(femaleEntries) > 0 {
		queenEntry = &femaleEntries[0]
	}
	if len(nonbinaryEntries) > 0 {
		nonbinaryEntry = &nonbinaryEntries[0]
	}
	if len(otherEntries) > 0 {
		otherEntry = &otherEntries[0]
	}

	// Get winner discord names for display purposes
	var kingsName string = ""
	var queensName string = ""
	var nonbinsName string = ""
	var othersName string = ""
	if kingEntry != nil {
		kingApiUser, err := s.User(kingEntry.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving king's API profile: %v", err)
		}
		kingsName = kingApiUser.Username
	}
	if queenEntry != nil {
		queenApiUser, err := s.User(queenEntry.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving queen's API profile: %v", err)
		}
		queensName = queenApiUser.Username
	}
	if nonbinaryEntry != nil {
		nonbinApiUser, err := s.User(nonbinaryEntry.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving nonbinary's API profile: %v", err)
		}
		nonbinsName = nonbinApiUser.Username
	}
	if otherEntry != nil {
		othersApiUser, err := s.User(otherEntry.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving other's API profile: %v", err)
		}
		othersName = othersApiUser.Username
	}

	now := time.Now()
	month := now.Format("January")
	year := now.Year()
	leaderboardMonthString := fmt.Sprintf("%s, %d", month, year)

	// Build winners embed
	embed := embed.
		NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖	Monthly Leaderboard Current State, `%s`", leaderboardMonthString)).
		SetDescription("The following OTA members have been the most active users this month (so far!) by engaging in conversations, receiving awards and spending time in voice channels.").
		SetColor(000000).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		AddLineBreakField()

	// If no valid entries found
	if kingsName == "" && queensName == "" && nonbinsName == "" && othersName == "" {
		embed.AddField("", "There are no valid monthly leaderboard entries at the moment.", false)
		return []*discordgo.MessageEmbed{embed.MessageEmbed}
	}

	if kingsName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(kingEntry.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("♂ Best so far, `%s`", kingsName), fieldValue, false)
	}

	if queensName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(queenEntry.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("♀ Best so far, `%s`", queensName), fieldValue, false)
	}

	if nonbinsName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(nonbinaryEntry.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("⚥ Best so far, `%s`", nonbinsName), fieldValue, false)
	}

	if othersName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(otherEntry.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("🌈 Best so far, `%s`", othersName), fieldValue, false)
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}

}
