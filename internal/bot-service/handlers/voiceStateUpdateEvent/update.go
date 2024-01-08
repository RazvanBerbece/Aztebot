package voiceStateUpdateEvent

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {

	var voiceChannels map[string]string
	if globals.Environment == "staging" {
		// Dev text channels
		voiceChannels = map[string]string{
			"1173790229258326106": "radio",
		}
	} else {
		// Production text channels
		voiceChannels = map[string]string{
			"1176204022399631381": "radio",
			"1118202946455351388": "music-1",
			"1118202975026937948": "music-2",
			"1118202999504904212": "music-3",
		}
		fmt.Println(voiceChannels)
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

	userId := member.User.ID

	if vs.SelfStream && vs.ChannelID != "" {
		// User STARTED STREAMING
		now := time.Now()
		globals.StreamSessions[userId] = &now

		err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
		}

		err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(userId, time.Now().Unix())
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", userId, err)
		}
	} else if vs.SelfStream && vs.ChannelID == "" {
		// The Discord API does something weird and sends SelfStream as true
		// when a user leaves a VC directly without stopping streaming first
		if joinTime, ok := globals.VoiceSessions[userId]; ok {
			duration := time.Since(joinTime)
			err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(duration.Seconds()))
			if err != nil {
				fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
			}
			delete(globals.VoiceSessions, userId)
			delete(globals.StreamSessions, userId)
			delete(globals.MusicSessions, userId)
		}
	} else {
		if vs.ChannelID != "" && globals.StreamSessions[userId] == nil && globals.MusicSessions[userId] == nil {
			// User JOINED a VC but NOT STREAMING
			if TargetChannelIsForMusicListening(voiceChannels, vs.ChannelID) {
				fmt.Println("JOIN MUSIC")
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
				fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", userId, err)
			}
		} else if vs.ChannelID != "" && globals.StreamSessions[userId] != nil {
			delete(globals.StreamSessions, userId)
		} else if vs.ChannelID == "" && globals.StreamSessions[userId] == nil {
			// User LEFT THE VOICE CHANNEL
			musicSession, userHadMusicSession := globals.MusicSessions[userId]
			if userHadMusicSession {
				// User was on a music channel
				fmt.Println("LEAVE MUSIC")
				for _, joinTime := range musicSession {
					duration := time.Since(*joinTime)
					err := globalsRepo.UserStatsRepository.AddToTimeSpentListeningMusic(userId, int(duration.Seconds()))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", userId, err)
					}
				}
			} else {
				// User was on any other VC
				if joinTime, ok := globals.VoiceSessions[userId]; ok {
					duration := time.Since(joinTime)
					err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(duration.Seconds()))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
					}
				}
			}
			delete(globals.MusicSessions, userId)
			delete(globals.VoiceSessions, userId)
			delete(globals.StreamSessions, userId)
		}
	}

}

func TargetChannelIsForMusicListening(voiceChannels map[string]string, channelId string) bool {
	fmt.Println("TARGET VC", channelId)
	for id := range voiceChannels {
		if channelId == id {
			// Target VC is a music-specific channel
			return true
		}
	}
	return false
}
