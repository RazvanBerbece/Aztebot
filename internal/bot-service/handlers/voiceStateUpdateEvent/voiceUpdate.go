package voiceStateUpdateEvent

import (
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {

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

	guild, err := s.Guild(vs.GuildID)
	if err != nil {
		log.Println("Error getting guild: ", err)
		return
	}

	member, err := s.GuildMember(guild.ID, vs.UserID)
	if err != nil {
		log.Println("Error getting member: ", err)
		return
	}

	if member.User.Bot {
		// Ignore bot voice state updates
		return
	}

	userId := member.User.ID

	if vs.SelfStream && vs.ChannelID != "" {
		// User STARTED STREAMING
		now := time.Now()
		globals.StreamSessions[userId] = &now

		err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
		}

		err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(userId, now.Unix())
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", userId, err)
		}
	} else if vs.SelfStream && vs.ChannelID == "" {
		// The Discord API does something weird and sends SelfStream as true
		// when a user leaves a VC directly without stopping streaming first
		if joinTime, ok := globals.VoiceSessions[userId]; ok {

			duration := time.Since(joinTime)
			secondsSpent := duration.Seconds()

			err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(secondsSpent))
			if err != nil {
				fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
			}

			// Publish experience grant message on the channel
			globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
				UserId:   userId,
				Points:   globals.ExperienceReward_InMusic * secondsSpent,
				Activity: "Time Spent Streaming",
			}

			delete(globals.VoiceSessions, userId)
			delete(globals.StreamSessions, userId)
			delete(globals.MusicSessions, userId)
		}
	} else {
		if vs.ChannelID != "" && globals.StreamSessions[userId] == nil && globals.MusicSessions[userId] == nil {
			// User JOINED a VC but NOT STREAMING
			if utils.TargetChannelIsForMusicListening(musicChannels, vs.ChannelID) {
				now := time.Now()
				globals.MusicSessions[userId] = map[string]*time.Time{
					vs.ChannelID: &now,
				}
			} else {
				globals.VoiceSessions[userId] = time.Now()
			}

			err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
			if err != nil {
				fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
			}

			err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(userId, time.Now().Unix())
			if err != nil {
				fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", userId, err)
			}
		} else if vs.ChannelID != "" && globals.StreamSessions[userId] != nil {
			delete(globals.StreamSessions, userId)
		} else if vs.ChannelID == "" && globals.StreamSessions[userId] == nil {
			// User LEFT THE VOICE CHANNEL
			musicSession, userHadMusicSession := globals.MusicSessions[userId]
			if userHadMusicSession {
				// User was on a music channel
				for _, joinTime := range musicSession {

					duration := time.Since(*joinTime)
					secondsSpent := duration.Seconds()

					err := globalsRepo.UserStatsRepository.AddToTimeSpentListeningMusic(userId, int(secondsSpent))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", userId, err)
					}

					// Publish experience grant message on the channel
					globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
						UserId:   userId,
						Points:   globals.ExperienceReward_InMusic * secondsSpent,
						Activity: "Time Spent in Music Channels",
					}

				}
			} else {
				// User was on any other VC
				if joinTime, ok := globals.VoiceSessions[userId]; ok {

					duration := time.Since(joinTime)
					secondsSpent := duration.Seconds()

					err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(duration.Seconds()))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
					}

					// Publish experience grant message on the channel
					globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
						UserId:   userId,
						Points:   globals.ExperienceReward_InVc * secondsSpent,
						Activity: "Time Spent in Voice Channels",
					}
				}
			}
			delete(globals.MusicSessions, userId)
			delete(globals.VoiceSessions, userId)
			delete(globals.StreamSessions, userId)
			delete(globals.DeafSessions, userId)
		}
	}

}

func UserHasActiveVoiceSession(uid string) bool {

	status := 0

	if _, ok := globals.VoiceSessions[uid]; ok {
		status += 1
	}

	if _, ok := globals.MusicSessions[uid]; ok {
		status += 1
	}

	if _, ok := globals.StreamSessions[uid]; ok {
		status += 1
	}

	return status == 3

}
