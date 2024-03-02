package readyEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/tasks/channelHandlers"
	cron "github.com/RazvanBerbece/Aztebot/internal/api/tasks/cron"
	"github.com/RazvanBerbece/Aztebot/internal/api/tasks/startup"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	utils.LogHandlerCall("Ready", "")

	// Load static data once runtime is confirmed
	go LoadStaticData()

	// Retrieve list of DB users at startup time (for convenience and some optimisation further down the line)
	uids, err := globalsRepo.UsersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Printf("Failed to load users at startup time: %v", err)
	}

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Other setups

	// Initial sync of members on server with the database
	go startup.SyncUsersAtStartup(s)

	// Initial cleanup of members from database against the Discord server
	go startup.CleanupMemberAtStartup(s, uids)

	// Initial update of experience gains in the DB
	go startup.SyncExperiencePointsGainsAtStartup(s)

	// Initial informative messages on certain channels
	go startup.SendInformationEmbedsToTextChannels(s)

	// Check for users on voice channels and start their VC sessions
	go startup.RegisterUsersInVoiceChannelsAtStartup(s)

	// Run background task to periodically update voice session durations in the DB
	go cron.UpdateVoiceSessionDurations(s)

	// Run channel message handlers
	go channelHandlers.HandleExperienceGrantsMessages(false)

	// CRON FUNCTIONS FOR VARIOUS FEATURES (like activity streaks, XP gaining?, etc.)
	cron.ProcessUpdateActivityStreaks(24, 0, 0)               // the hh:mm:ss timestamp in a day to run the cron at
	cron.ProcessRemoveExpiredWarns(2)                         // run every n=2 months
	cron.ClearExpiredTimeouts(s)                              // clear timeouts with freq from env var
	cron.ProcessRemoveArchivedTimeouts(1)                     // run every n=1 month
	cron.ProcessMonthlyLeaderboard(s, 23, 55, 00, true, true) // run on last day at given time

}
