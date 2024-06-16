package readyEvent

import (
	"fmt"
	"time"

	cron "github.com/RazvanBerbece/Aztebot/internal/bot-service/api/tasks/cron"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/tasks/startup"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	logging.LogHandlerCall("Ready", "")

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

	// Initial informative messages on certain channels
	go startup.SendInformationEmbedsToTextChannels(s)

	// Check for users on voice channels and start their VC sessions
	go RegisterUsersInVoiceChannelsAtStartup(s)

	// Run background task to periodically update voice session durations in the DB
	go UpdateVoiceSessionDurations(s)

	// CRON FUNCTIONS FOR VARIOUS FEATURES (like activity streaks, XP gaining?, etc.)
	cron.ProcessUpdateActivityStreaks(24, 0, 0) // the hh:mm:ss timestamp in a day to run the cron at
	cron.ProcessRemoveExpiredWarns(2)           // run every n=2 months

}

func UpdateVoiceSessionDurations(s *discordgo.Session) {

	var numSec int
	if globals.UpdateVoiceStateFrequencyErr != nil {
		numSec = 5
	} else {
		numSec = globals.UpdateVoiceStateFrequency
	}

	fmt.Println("Starting Task UpdateVoiceSesionDurations() at", time.Now(), "running every", numSec, "seconds")

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go updateVoiceSessions()
				go updateStreamingSessions()
				go updateMusicSessions()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func RegisterUsersInVoiceChannelsAtStartup(s *discordgo.Session) {

	fmt.Println("Trying to RegisterUsersInVoiceChannelsAtStartup() at", time.Now())

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

			user, err := s.User(userId)
			if err != nil {
				fmt.Println("Error retrieving user:", err)
				return
			}
			if user.Bot {
				continue
			}

			if utils.TargetChannelIsForMusicListening(musicChannels, channelId) {
				// If the voice state is purposed for music, initiate a music session at startup time
				_, exists := globals.MusicSessions[userId]
				if exists {
					continue
				} else {
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
						globals.StreamSessions[userId] = &now
						streamSessionsAtStartup += 1
					}
				} else {
					// If the voice state is purposed for just for listening on a voice channel, initiate a voice session at startup time
					_, exists := globals.VoiceSessions[userId]
					if exists {
						continue
					} else {
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
		fmt.Printf("Found %d active voice states at bot startup time (%d voice, %d streaming, %d music)\n", totalSessions, voiceSessionsAtStartup, streamSessionsAtStartup, musicSessionsAtStartup)
	}

}

func updateVoiceSessions() {
	for uid, joinTime := range globals.VoiceSessions {
		duration := time.Since(joinTime)
		err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globals.VoiceSessions[uid] = now
	}
}

func updateStreamingSessions() {
	for uid, joinTime := range globals.StreamSessions {
		duration := time.Since(*joinTime)
		err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding streaming duration to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globals.StreamSessions[uid] = &now
	}
}

func updateMusicSessions() {
	for uid := range globals.MusicSessions {
		session, userHadMusicSession := globals.MusicSessions[uid]
		if userHadMusicSession {
			// User was on a music channel
			for channelId, joinTime := range session {
				duration := time.Since(*joinTime)
				err := globalsRepo.UserStatsRepository.AddToTimeSpentListeningMusic(uid, int(duration.Seconds()))
				if err != nil {
					fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", uid, err)
				}

				// Reset join time
				now := time.Now()
				globals.MusicSessions[uid] = map[string]*time.Time{
					channelId: &now,
				}
			}
		}
	}
}
