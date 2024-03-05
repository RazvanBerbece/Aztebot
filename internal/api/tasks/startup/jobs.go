package startup

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncUsersAtStartup(s *discordgo.Session) error {

	fmt.Println("[STARTUP] Starting Task SyncUsersAtStartup() at", time.Now())

	// Inject new connections
	rolesRepository := repositories.NewRolesRepository()
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globals.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("[STARTUP] Failed Task SyncUsersAtStartup() at", time.Now(), "with error", err)
		return err
	}

	// Process the current batch of members
	processMembers(s, members, rolesRepository, usersRepository, userStatsRepository)

	// Paginate
	for len(members) == 1000 {
		// Set the 'After' parameter to the ID of the last member in the current batch
		lastMemberID := members[len(members)-1].User.ID
		members, err = s.GuildMembers(globals.DiscordMainGuildId, lastMemberID, 1000)
		if err != nil {
			fmt.Println("[STARTUP] Failed Task SyncUsersAtStartup() at", time.Now(), "with error", err)
			return err
		}

		// Process the next batch of members
		processMembers(s, members, rolesRepository, usersRepository, userStatsRepository)
	}

	// Cleanup
	go utils.CleanupRepositories(rolesRepository, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task SyncUsersAtStartup() at", time.Now())

	return nil

}

func CleanupMemberAtStartup(s *discordgo.Session, uids []string) error {

	fmt.Println("[STARTUP] Starting Task CleanupMemberAtStartup() at", time.Now())

	// Inject new connections
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	uidsLength := len(uids)

	// For each tag in the DB, delete user from table
	for i := 0; i < uidsLength; i++ {
		uid := uids[i]
		_, err := s.GuildMember(globals.DiscordMainGuildId, uid)
		if err != nil {
			// if the member does not exist on the main server, delete from the database
			err = member.DeleteAllMemberData(uid)
			if err != nil {
				fmt.Println("Error deleting hanging user data on startup sync: ", err)
			}
		}
	}

	// Cleanup repos
	go utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task CleanupMemberAtStartup() at", time.Now())

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) {
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := utils.SyncUserPersistent(s, globals.DiscordMainGuildId, member.User.ID, member, rolesRepository, usersRepository, userStatsRepository)
		if err != nil && err.Error() != "no update was executed" {
			fmt.Printf("Error syncing member %s: %v\n", member.User.Username, err)
		}
	}
}

func SendInformationEmbedsToTextChannels(s *discordgo.Session) {

	var textChannels map[string]string

	// TODO: Make the channels and their descriptions use environment variables somehow
	if globals.Environment == "staging" {
		// Dev text channels
		textChannels = map[string]string{
			"1188135110042734613": "default",
			"1194451477192773773": "staff-rules",
			"1198686819928264784": "server-rules",
			"1205859615406030868": "legends",
		}
	} else {
		// Production text channels
		textChannels = map[string]string{
			"1176277764001767464": "info-music",
			"1100486860058398770": "staff-rules",
			"1100142572141281460": "server-rules",
			"1100762035450544219": "legends",
		}
	}

	// For each available default message resource in local storage
	for id, details := range textChannels {
		hasMessage, err := utils.ChannelHasDefaultInformationMessage(s, id)
		if err != nil {
			fmt.Printf("Could not check for default message in channel %s (%s): %v", id, details, err)
			continue
		}
		if hasMessage {
			// Do not send this default message as it already exists
			continue
		} else {
			// Send associated default message to given text channel
			var embedText string
			var hasOwnEmbed bool
			var longEmbed *embed.Embed
			switch details {
			case "default":
				embedText = utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/default.txt")
			case "info-music":
				embedText = utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/music-info.txt")
			case "staff-rules":
				embedText = utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/staff-rules.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			case "server-rules":
				embedText = utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/server-rules.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			case "legends":
				embedText = utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/legends.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			}

			var messageEmbedToPost *discordgo.MessageEmbed
			if !hasOwnEmbed {
				messageEmbedToPost = embed.NewEmbed().
					SetTitle("🤖  Information Message").
					SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
					SetColor(000000).
					AddField("", embedText, false).
					MessageEmbed
			} else {
				messageEmbedToPost = longEmbed.MessageEmbed
			}

			_, err := s.ChannelMessageSendEmbed(id, messageEmbedToPost)
			if err != nil {
				log.Fatalf("An error occured while sending a default message (%s): %v", details, err)
				return
			}
		}
	}

}

func RegisterUsersInVoiceChannelsAtStartup(s *discordgo.Session) {

	fmt.Println("[STARTUP] Starting Task RegisterUsersInVoiceChannelsAtStartup() at", time.Now())

	now := time.Now()

	var musicChannels map[string]string
	if globals.Environment == "staging" {
		// Dev text channels
		musicChannels = map[string]string{
			"1173790229258326106": "radio",
		}
	} else {
		// Production text channels
		musicChannels = map[string]string{
			"1176204022399631381": "radio",
			"1118202946455351388": "music-1",
			"1118202975026937948": "music-2",
			"1118202999504904212": "music-3",
		}
	}

	guild, err := s.State.Guild(globals.DiscordMainGuildId)
	if err != nil {
		fmt.Println("Error retrieving guild:", err)
		return
	}

	// For each active voice state in the guild
	var voiceSessionsAtStartup int = 0
	var streamSessionsAtStartup int = 0
	var musicSessionsAtStartup int = 0
	var deafSessionsAtStartup int = 0

	var loadedUsersFromVCs bool = false
	var loadingTimeIsUp bool = false

	for !loadedUsersFromVCs {

		time.Sleep(5 * time.Millisecond)

		durationForLoadingSessions := time.Since(now)
		if durationForLoadingSessions.Seconds() > 2*60 { // only try this for ~2-3 minutes, then break and return
			loadingTimeIsUp = true
			break
		}

		for _, voiceState := range guild.VoiceStates {

			userId := voiceState.UserID
			channelId := voiceState.ChannelID

			userIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, userId, false)
			if err != nil {
				fmt.Println("Error retrieving user for bot check:", err)
				return
			}
			if *userIsBot {
				continue
			}

			if utils.TargetChannelIsForMusicListening(musicChannels, channelId) {
				// If the voice state is purposed for music, initiate a music session at startup time
				_, exists := globals.MusicSessions[userId]
				if exists {
					continue
				} else {
					now = time.Now()
					globals.MusicSessions[userId] = map[string]*time.Time{
						channelId: &now,
					}
					musicSessionsAtStartup += 1
				}
			} else {
				if voiceState.SelfStream {
					// If the voice state is purposed for streaming, initiate a streaming session at startup time
					_, exists := globals.StreamSessions[userId]
					if exists {
						continue
					} else {
						now = time.Now()
						globals.StreamSessions[userId] = &now
						streamSessionsAtStartup += 1
					}
				} else {
					// If the voice state is purposed for just for listening on a voice channel, initiate a voice session at startup time
					_, exists := globals.VoiceSessions[userId]
					if exists {
						continue
					} else {
						now = time.Now()
						globals.VoiceSessions[userId] = now
						voiceSessionsAtStartup += 1
					}
				}
			}

			loadedUsersFromVCs = true
		}

	}

	if loadedUsersFromVCs || loadingTimeIsUp {
		totalSessions := voiceSessionsAtStartup + streamSessionsAtStartup + musicSessionsAtStartup
		fmt.Printf("[STARTUP] Found %d active voice states at bot startup time (%d voice, %d streaming, %d music, %d deafened)\n", totalSessions, voiceSessionsAtStartup, streamSessionsAtStartup, musicSessionsAtStartup, deafSessionsAtStartup)
	}

}

func SyncExperiencePointsGainsAtStartup(s *discordgo.Session) {

	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "with error", err)
	}

	// For all users in the database
	fmt.Println("[STARTUP] Checkpoint Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "-> Updating", len(uids), "XP gains")
	for _, uid := range uids {

		user, err := usersRepository.GetUser(uid)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", err)
		}

		stats, errStats := userStatsRepository.GetStatsForUser(uid)
		if errStats != nil {
			if errStats == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", errStats)
		}

		updatedXp := utils.CalculateExperiencePointsFromStats(
			stats.NumberMessagesSent,
			stats.NumberSlashCommandsUsed,
			stats.NumberReactionsReceived,
			stats.TimeSpentInVoiceChannels,
			stats.TimeSpentListeningToMusic,
			globals.ExperienceReward_MessageSent,
			globals.ExperienceReward_SlashCommandUsed,
			globals.ExperienceReward_ReactionReceived,
			globals.ExperienceReward_InVc,
			globals.ExperienceReward_InMusic)
		user.CurrentExperience = float64(updatedXp)

		// Update user entity with new XP value
		_, err = usersRepository.UpdateUser(*user)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", err)
		}

	}

	// Cleanup repos
	go utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task SyncExperiencePointsGainsAtStartup() at", time.Now())

}