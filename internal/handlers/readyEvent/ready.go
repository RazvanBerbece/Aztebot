package readyEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/tasks/channelHandlers"
	cron "github.com/RazvanBerbece/Aztebot/internal/api/tasks/cron"
	"github.com/RazvanBerbece/Aztebot/internal/api/tasks/startup"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	utils.LogHandlerCall("Ready", "")

	// Load static data once Discord API runtime features are confirmed
	LoadStaticData()

	// Retrieve list of DB users at startup time (for convenience and some optimisation further down the line)
	uids, err := globalsRepo.UsersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Printf("Failed to load users at startup time: %v", err)
	}

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Initial sync of members on server with the database
	go startup.SyncUsersAtStartup(s)

	// Initial cleanup of members from database against the Discord server
	go startup.CleanupMemberAtStartup(s, uids)

	// Initial update of experience gains in the DB
	go startup.SyncExperiencePointsGainsAtStartup(s)

	// Initial publishing of informative messages on certain channels
	go startup.SendInformationEmbedsToTextChannels(s)

	// Check for users on voice channels and start their VC sessions
	go startup.RegisterUsersInVoiceChannelsAtStartup(s)

	// Run background task to periodically update voice session durations in the DB
	go cron.UpdateVoiceSessionDurations(s)

	// Run event handlers
	go channelHandlers.HandleNotificationEvents(s)
	go channelHandlers.HandleExperienceGrantEvents()
	go channelHandlers.HandleDynamicChannelCreationEvents(s)

	// CRON FUNCTIONS FOR VARIOUS FEATURES (like activity streaks, cleanups, etc.)
	cron.ProcessUpdateActivityStreaks(24, 0, 0)               // the hh:mm:ss timestamp in a day to run the cron at (i.e 24:00:00)
	cron.ProcessMonthlyLeaderboard(s, 23, 55, 0, true, false) // run on last day of current month at given time (i.e 23:55:00)
	cron.ProcessClearExpiredTimeouts(s)
	cron.ProcessCleanupUnusedDynamicChannels(s, globals.DiscordMainGuildId)
	cron.ProcessRemoveExpiredWarns()
	cron.ProcessRemoveArchivedTimeouts()

}
