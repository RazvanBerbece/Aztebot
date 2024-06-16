package cron

import (
	"fmt"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
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
				go updateVoiceSessions(userStatsRepository)
				go updateStreamingSessions(userStatsRepository)
				go updateMusicSessions(userStatsRepository)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func updateVoiceSessions(userStatsRepo *repositories.UsersStatsRepository) {
	for uid, joinTime := range globals.VoiceSessions {

		duration := time.Since(joinTime)
		secondsSpent := duration.Seconds()

		err := userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globals.VoiceSessions[uid] = now

		// Publish experience grant message on the channel
		globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
			UserId:   uid,
			Points:   globals.ExperienceReward_InMusic * secondsSpent,
			Activity: "Time Spent in Voice Channels",
		}
	}
}

func updateStreamingSessions(userStatsRepo *repositories.UsersStatsRepository) {
	for uid, joinTime := range globals.StreamSessions {

		duration := time.Since(*joinTime)
		secondsSpent := duration.Seconds()

		err := userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(secondsSpent))
		if err != nil {
			fmt.Printf("An error ocurred while adding streaming duration to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globals.StreamSessions[uid] = &now

		// Publish experience grant message on the channel
		globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
			UserId:   uid,
			Points:   globals.ExperienceReward_InMusic * secondsSpent,
			Activity: "Time Spent Streaming",
		}
	}
}

func updateMusicSessions(userStatsRepo *repositories.UsersStatsRepository) {
	for uid := range globals.MusicSessions {

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

				// Reset join time
				now := time.Now()
				globals.MusicSessions[uid] = map[string]*time.Time{
					channelId: &now,
				}

				// Publish experience grant message on the channel
				globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
					UserId:   uid,
					Points:   globals.ExperienceReward_InMusic * secondsSpent,
					Activity: "Time Spent in Music Channels",
				}
			}
		}
	}
}
