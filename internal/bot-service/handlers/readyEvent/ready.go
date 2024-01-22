package readyEvent

import (
	"fmt"
	"log"
	"strings"
	"time"

	cronFeature "github.com/RazvanBerbece/Aztebot/internal/bot-service/api/cron/feature"
	cronUser "github.com/RazvanBerbece/Aztebot/internal/bot-service/api/cron/user"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
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
	go cronUser.SyncUsersAtStartup(s)

	// Initial cleanup of members from database against the Discord server
	go cronUser.CleanupMemberAtStartup(s, uids)

	// Initial informative messages on certain channels
	go SendInformationEmbedsToTextChannels(s)

	// Check for users on voice channels and start their VC sessions
	go RegisterUsersInVoiceChannelsAtStartup(s)

	// Run background task toperiodically update voice session durations in the DB
	go UpdateVoiceSessionDurations(s)

	// CRON FUNCTIONS FOR VARIOUS FEATURES (like activity streaks, XP gaining?, etc.)
	cronFeature.ProcessUpdateActivityStreaks(24, 0, 0) // the hh:mm:ss timestamp in a day to run the cron at
	cronFeature.ProcessRemoveExpiredWarns(2)           // run every n=2 months

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
		}
	} else {
		// Production text channels
		textChannels = map[string]string{
			"1176277764001767464": "info-music",
			"1100486860058398770": "staff-rules",
			"1100142572141281460": "server-rules",
		}
	}

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
			var ownEmbed *embed.Embed = embed.NewEmbed()
			switch details {
			case "default":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/default.txt")
			case "info-music":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/music-info.txt")
			case "staff-rules":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/staff-rules.txt")
				hasOwnEmbed = true
				mutateLongEmbedFromStaticData(embedText, ownEmbed)
			case "server-rules":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/server-rules.txt")
				hasOwnEmbed = true
				mutateLongEmbedFromStaticData(embedText, ownEmbed)
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
				messageEmbedToPost = ownEmbed.MessageEmbed
			}

			_, err := s.ChannelMessageSendEmbed(id, messageEmbedToPost)
			if err != nil {
				log.Fatalf("An error occured while sending a default message (%s): %v", details, err)
				return
			}
		}
	}

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

// Note that this is a mutating function on `hasOwnEmbed` and `embed`.
func mutateLongEmbedFromStaticData(embedText string, embed *embed.Embed) {
	// Split the content into sections based on double newline characters ("\n\n")
	sections := strings.Split(embedText, "\n\n")
	for _, section := range sections {
		lines := strings.Split(section, "\n")
		if len(lines) > 0 {
			// Use the first line as the title and the rest as content
			title := lines[0]
			content := strings.Join(lines[1:], "\n")
			embed.AddField(title, content, false)
		}
	}
}
