package cron

import (
	"fmt"
	"time"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

// Process the monthly leaderboard results at the given h:m:s timestamp.
// actualLastDay defines whether this will execute on the actual day of the month. if false, it will execute on current day.
// dryrun defines whether this will clear out the monthlyLeaderboard table after execution. if false, it will leave the table in place.
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

	var kingEntry *dax.MonthlyLeaderboardEntry = nil
	var queenEntry *dax.MonthlyLeaderboardEntry = nil
	var nonbinaryEntry *dax.MonthlyLeaderboardEntry = nil
	var otherEntry *dax.MonthlyLeaderboardEntry = nil
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
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-globalAnnouncements"]; channelExists {
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

func sendMonthlyLeaderboardWinnerNotification(s *discordgo.Session, channelId string, king *dax.MonthlyLeaderboardEntry, queen *dax.MonthlyLeaderboardEntry, nonbinary *dax.MonthlyLeaderboardEntry, other *dax.MonthlyLeaderboardEntry) {

	// Get winner discord names for display purposes
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
		SetDescription("The following OTA members have been the most active users this month by engaging in conversations, receiving awards and spending time in voice channels.").
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

	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           embed,
	}

}
