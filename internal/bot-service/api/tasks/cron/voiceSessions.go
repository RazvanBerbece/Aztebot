package cron

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

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

func updateVoiceSessions() {
	for uid, joinTime := range globals.VoiceSessions {

		// if _, isDeafened := globals.DeafSessions[uid]; isDeafened {
		// 	fmt.Println("SKIP VS")
		// 	// Skip adding time for deafened users
		// 	// But reset the join time
		// 	now := time.Now()
		// 	globals.VoiceSessions[uid] = now
		// 	continue
		// }

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

		// if _, isDeafened := globals.DeafSessions[uid]; isDeafened {
		// 	// Skip adding time for deafened users
		// 	// But reset the join time
		// 	fmt.Println("SKIP STREAM")
		// 	now := time.Now()
		// 	globals.StreamSessions[uid] = &now
		// 	continue
		// }

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

				// if _, isDeafened := globals.DeafSessions[uid]; isDeafened {
				// 	fmt.Println("SKIP MUSIC")
				// 	// Skip adding time for deafened users
				// 	// But reset the join time
				// 	now := time.Now()
				// 	globals.MusicSessions[uid] = map[string]*time.Time{
				// 		channelId: &now,
				// 	}
				// 	continue
				// }

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
