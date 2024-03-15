package cron

import (
	"fmt"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func ProcessMonthlyLeaderboard(s *discordgo.Session, hour int, minute int, second int, actualLastDay bool, dryrun bool) {

	initialMonthlyLeaderboardDelay, monthlyLeaderboardTicker := GetDelayAndTickerForMonthlyLeaderboardCron(actualLastDay, hour, minute, second)

	go func() {

		hoursAsDays := initialMonthlyLeaderboardDelay.Hours() / 24
		fmt.Println("[SCHEDULED CRON] Scheduled Task ExtractMonthlyLeaderboardWinners() in <", hoursAsDays, "> days")
		time.Sleep(initialMonthlyLeaderboardDelay)

		// Inject new connections
		monthlyLeaderboardRepository := repositories.NewMonthlyLeaderboardRepository()

		// The first run should happen at start-up, not after n days
		ExtractMonthlyLeaderboardWinners(s, monthlyLeaderboardRepository, dryrun)

		for range monthlyLeaderboardTicker.C {
			// Process
			ExtractMonthlyLeaderboardWinners(s, monthlyLeaderboardRepository, dryrun)
		}
	}()
}

func ExtractMonthlyLeaderboardWinners(s *discordgo.Session, monthlyLeaderboardRepository *repositories.MonthlyLeaderboardRepository, dryrun bool) {

	fmt.Println("[CRON] Starting Task ExtractMonthlyLeaderboardWinners() at", time.Now())

	// Extract each category winner
	// TODO: Could make repository function return only the winner - if decided that other entries don't matter in the end
	maleEntries, err := monthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(0)
	if err != nil {
		fmt.Println("[CRON] Failed Task ExtractMonthlyLeaderboardWinners() at", time.Now(), "with error", err)
	}
	femaleEntries, err := monthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(1)
	if err != nil {
		fmt.Println("[CRON] Failed Task ExtractMonthlyLeaderboardWinners() at", time.Now(), "with error", err)
	}
	nonbinaryEntries, err := monthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(2)
	if err != nil {
		fmt.Println("[CRON] Failed Task ExtractMonthlyLeaderboardWinners() at", time.Now(), "with error", err)
	}
	otherEntries, err := monthlyLeaderboardRepository.GetLeaderboardEntriesByCategory(3)
	if err != nil {
		fmt.Println("[CRON] Failed Task ExtractMonthlyLeaderboardWinners() at", time.Now(), "with error", err)
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

	// Send winner notification to designated channel
	if channel, channelExists := globals.NotificationChannels["notif-monthlyWinners"]; channelExists {
		go sendMonthlyLeaderboardWinnerNotification(s, channel.ChannelId, kingEntry, queenEntry, nonbinaryEntry, otherEntry)
	}

	// Reset the leaderboard so it can be used for the next month (if not a dryrun)
	if !dryrun {
		err = monthlyLeaderboardRepository.ClearLeaderboard()
		if err != nil {
			fmt.Println("[CRON] Failed Task ExtractMonthlyLeaderboardWinners() at", time.Now(), "with error", err)
		}
	}

	fmt.Println("[CRON] Finished Task ExtractMonthlyLeaderboardWinners() at", time.Now())

}

func sendMonthlyLeaderboardWinnerNotification(s *discordgo.Session, channelId string, king *dataModels.MonthlyLeaderboardEntry, queen *dataModels.MonthlyLeaderboardEntry, nonbinary *dataModels.MonthlyLeaderboardEntry, other *dataModels.MonthlyLeaderboardEntry) {

	// Get command owner discord names
	var kingsName string = ""
	var queensName string = ""
	var nonbinsName string = ""
	var othersName string = ""
	if king != nil {
		kingApiUser, err := s.User(king.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving king's API profile: %v", err)
		}
		kingsName = kingApiUser.Username
	}
	if queen != nil {
		queenApiUser, err := s.User(queen.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving queen's API profile: %v", err)
		}
		queensName = queenApiUser.Username
	}
	if nonbinary != nil {
		nonbinApiUser, err := s.User(nonbinary.UserId)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving nonbinary's API profile: %v", err)
		}
		nonbinsName = nonbinApiUser.Username
	}
	if other != nil {
		othersApiUser, err := s.User(other.UserId)
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
		SetTitle(fmt.Sprintf("🤖📣 Monthly Leaderboard Winners Announcement, `%s` 📣📣", leaderboardMonthString)).
		SetColor(000000).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		AddLineBreakField()

	if kingsName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(king.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("♂ King of The Month, `%s`", kingsName), fieldValue, false)
	}

	if queensName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(queen.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("♀ Queen of The Month, `%s`", queensName), fieldValue, false)
	}

	if nonbinsName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(nonbinary.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("⚥ Nonbinary of The Month, `%s`", nonbinsName), fieldValue, false)
	}

	if othersName != "" {
		fieldValue := fmt.Sprintf("Accumulated a total of 💠 `%d` XP !", int64(other.XpEarnedInCurrentMonth))
		embed.AddField(fmt.Sprintf("🌈 Others of The Month, `%s`", othersName), fieldValue, false)
	}

	// Tag everyone to propagate announcement
	embed.
		AddLineBreakField().
		AtTagEveryone()

	// Publish notification event
	globals.NotificationsChannel <- events.NotificationEvent{
		Session:         s,
		TargetChannelId: channelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           embed,
	}

}
