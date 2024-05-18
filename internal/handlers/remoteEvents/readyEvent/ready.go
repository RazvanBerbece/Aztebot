package readyEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	channelHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/channelEvents"
	cron "github.com/RazvanBerbece/Aztebot/internal/tasks/cron"
	"github.com/RazvanBerbece/Aztebot/internal/tasks/startup"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	// Load static data once Discord API runtime features are confirmed
	LoadStaticData()

	// Retrieve list of DB users at startup time (for convenience and some optimisation further down the line)
	uids, err := globalRepositories.UsersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Printf("Failed to load users at startup time: %v", err)
	}

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Run gochannel event handlers
	go channelHandlers.HandleNotificationEvents(s)
	go channelHandlers.HandleExperienceGrantEvents()
	go channelHandlers.HandleDynamicChannelCreationEvents(s)
	go channelHandlers.HandleMemberMessageDeletionEvents(s)
	go channelHandlers.HandleDirectMessageEvents(s)
	go channelHandlers.HandleComplexResponseEvents(s, globalConfiguration.EmbedPageSize)
	go channelHandlers.HandlePromotionRequestEvents(s, globalConfiguration.OrderRoleNames, true, false)

	// Initial sync of members on server with the database
	go startup.SyncMembersAtStartup(s, globalConfiguration.OrderRoleNames, false)

	// Initial cleanup of members from database against the Discord server
	go startup.CleanupMemberAtStartup(s, uids)

	// Initial update of experience gains and current levels in the DB
	go startup.SyncExperiencePointsGainsAtStartup(s, uids)

	// Initial publishing of informative messages on certain channels
	go startup.SendInformationEmbedsToTextChannels(s)

	// Check for users on voice channels and start their VC sessions
	go startup.RegisterUsersInVoiceChannelsAtStartup(s)

	// Run background task to periodically update voice session durations in the DB
	go cron.UpdateVoiceSessionDurations(s)

	// CRON FEATS
	cron.ProcessUpdateActivityStreaks(24, 0, 0)               // the hh:mm:ss timestamp in a day to run the cron at (i.e 24:00:00)
	cron.ProcessMonthlyLeaderboard(s, 23, 55, 0, true, false) // run on last day of current month at given time (i.e 23:55:00)

	// CRON RUNTIME & PERSISTENT ENTITY CLEANUPS
	cron.ProcessClearExpiredTimeouts(s)
	cron.ProcessCleanupUnusedDynamicChannels(s, globalConfiguration.DiscordMainGuildId)
	cron.ProcessRemoveExpiredWarns()
	cron.ProcessRemoveArchivedTimeouts()

	// CRON RUNTIME STATE CLEANUPS
	cron.ClearOldPaginatedEmbeds(s)

}
