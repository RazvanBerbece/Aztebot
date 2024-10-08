package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/bwmarrin/discordgo"
)

func UpdateVoiceSessionDurations(s *discordgo.Session) {

	var numSec int
	if globalConfiguration.UpdateVoiceStateFrequencyErr != nil {
		numSec = 5
	} else {
		numSec = globalConfiguration.UpdateVoiceStateFrequency
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
	for uid, joinTime := range globalState.VoiceSessions {

		duration := time.Since(joinTime)
		secondsSpent := duration.Seconds()

		err := userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globalState.VoiceSessions[uid] = now

		// Grant XP
		globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
			UserId: uid,
			Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
			Type:   "VOICE_ACTIVITY",
		}

		// Award coins
		globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
			GuildId:  globalConfiguration.DiscordMainGuildId,
			UserId:   uid,
			Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
			Activity: "TIME-VC",
		}
	}
}

func updateStreamingSessions(userStatsRepo *repositories.UsersStatsRepository) {
	for uid, joinTime := range globalState.StreamSessions {

		duration := time.Since(*joinTime)
		secondsSpent := duration.Seconds()

		err := userStatsRepo.AddToTimeSpentInVoiceChannels(uid, int(secondsSpent))
		if err != nil {
			fmt.Printf("An error ocurred while adding streaming duration to voice channels for user with id %s: %v", uid, err)
		}

		// Reset join time
		now := time.Now()
		globalState.StreamSessions[uid] = &now

		// Grant XP
		globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
			UserId: uid,
			Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
			Type:   "VOICE_ACTIVITY",
		}

		// Award coins
		globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
			GuildId:  globalConfiguration.DiscordMainGuildId,
			UserId:   uid,
			Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
			Activity: "TIME-VC",
		}
	}
}

func updateMusicSessions(userStatsRepo *repositories.UsersStatsRepository) {
	for uid := range globalState.MusicSessions {

		session, userHadMusicSession := globalState.MusicSessions[uid]
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
				globalState.MusicSessions[uid] = map[string]*time.Time{
					channelId: &now,
				}

				// Grant XP
				globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
					UserId: uid,
					Points: globalConfiguration.ExperienceReward_InMusic * secondsSpent,
					Type:   "MUSIC_ACTIVITY",
				}

				// Award coins
				globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
					GuildId:  globalConfiguration.DiscordMainGuildId,
					UserId:   uid,
					Funds:    globalConfiguration.CoinReward_InMusic * secondsSpent,
					Activity: "TIME-MUSIC",
				}
			}
		}
	}
}
