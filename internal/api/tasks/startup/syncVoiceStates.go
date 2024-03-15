package startup

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterUsersInVoiceChannelsAtStartup(s *discordgo.Session) {

	fmt.Println("[STARTUP] Starting Task RegisterUsersInVoiceChannelsAtStartup() at", time.Now())

	now := time.Now()

	guild, err := s.State.Guild(globalConfiguration.DiscordMainGuildId)
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

			userIsBot, err := member.IsBot(s, globalConfiguration.DiscordMainGuildId, userId, false)
			if err != nil {
				fmt.Println("Error retrieving user for bot check:", err)
				return
			}
			if *userIsBot {
				continue
			}

			if utils.TargetChannelIsForMusicListening(globalConfiguration.MusicChannels, channelId) {
				// If the voice state is purposed for music, initiate a music session at startup time
				_, exists := globalState.MusicSessions[userId]
				if exists {
					continue
				} else {
					now = time.Now()
					globalState.MusicSessions[userId] = map[string]*time.Time{
						channelId: &now,
					}
					musicSessionsAtStartup += 1
				}
			} else {
				if voiceState.SelfStream {
					// If the voice state is purposed for streaming, initiate a streaming session at startup time
					_, exists := globalState.StreamSessions[userId]
					if exists {
						continue
					} else {
						now = time.Now()
						globalState.StreamSessions[userId] = &now
						streamSessionsAtStartup += 1
					}
				} else {
					// If the voice state is purposed for just for listening on a voice channel, initiate a voice session at startup time
					_, exists := globalState.VoiceSessions[userId]
					if exists {
						continue
					} else {
						now = time.Now()
						globalState.VoiceSessions[userId] = now
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
