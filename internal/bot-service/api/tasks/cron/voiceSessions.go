package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/bwmarrin/discordgo"
)

func UpdateVoiceSessionDurations(s *discordgo.Session) {

	var numSec int
	if globals.UpdateVoiceStateFrequencyErr != nil {
		numSec = 5
	} else {
		numSec = globals.UpdateVoiceStateFrequency
	}

	fmt.Println("[CRON] Starting Cron Ticker UpdateVoiceSesionDurations() at", time.Now(), "running every", numSec, "seconds")

	// Inject repositories
	userStatsRepository := repositories.NewUsersStatsRepository()

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go updateVoiceSessions(s, userStatsRepository)
				go updateStreamingSessions(s, userStatsRepository)
				go updateMusicSessions(s, userStatsRepository)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func updateVoiceSessions(s *discordgo.Session, userStatsRepo *repositories.UsersStatsRepository) {
	for uid, joinTime := range globals.VoiceSessions {

		// Ignore all bot sessions
		authorIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, uid, false)
		if err != nil {
			continue
		}
		if *authorIsBot {
			continue
		}

		duration := time.Since(joinTime)
		secondsSpent := duration.Seconds()

		err = userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", uid, err)
		}

		// Grant experience points for time spent streaming
		go member.GrantMemberExperience(uid, "IN_VC_REWARD", &secondsSpent)

		// Reset join time
		now := time.Now()
		globals.VoiceSessions[uid] = now
	}
}

func updateStreamingSessions(s *discordgo.Session, userStatsRepo *repositories.UsersStatsRepository) {
	for uid, joinTime := range globals.StreamSessions {

		// Ignore all bot sessions
		authorIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, uid, false)
		if err != nil {
			continue
		}
		if *authorIsBot {
			continue
		}

		duration := time.Since(*joinTime)
		secondsSpent := duration.Seconds()

		err = userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(secondsSpent))
		if err != nil {
			fmt.Printf("An error ocurred while adding streaming duration to voice channels for user with id %s: %v", uid, err)
		}

		// Grant experience points for time spent streaming
		go member.GrantMemberExperience(uid, "IN_VC_REWARD", &secondsSpent)

		// Reset join time
		now := time.Now()
		globals.StreamSessions[uid] = &now
	}
}

func updateMusicSessions(s *discordgo.Session, userStatsRepo *repositories.UsersStatsRepository) {
	for uid := range globals.MusicSessions {

		// Ignore all bot sessions
		authorIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, uid, false)
		if err != nil {
			continue
		}
		if *authorIsBot {
			continue
		}

		session, userHadMusicSession := globals.MusicSessions[uid]
		if userHadMusicSession {
			// User was on a music channel
			for channelId, joinTime := range session {

				duration := time.Since(*joinTime)
				secondsSpent := duration.Seconds()

				err := userStatsRepo.AddToTimeSpentListeningMusic(uid, int(secondsSpent))
				if err != nil {
					fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", uid, err)
				}

				// Grant experience points for time spent listening to music
				go member.GrantMemberExperience(uid, "IN_MUSIC_REWARD", &secondsSpent)

				// Reset join time
				now := time.Now()
				globals.MusicSessions[uid] = map[string]*time.Time{
					channelId: &now,
				}
			}
		}
	}
}
