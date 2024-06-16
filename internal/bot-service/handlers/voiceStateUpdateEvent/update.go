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

	if vs.SelfStream {
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
	} else {
		if vs.ChannelID != "" && globals.StreamSessions[userId] == nil {
			// User JOINED a VC but NOT STREAMING
			globals.VoiceSessions[userId] = time.Now()

			err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
			if err != nil {
				fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
			}

			err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(userId, time.Now().Unix())
			if err != nil {
				fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", userId, err)
			}
		} else if vs.ChannelID != "" && globals.StreamSessions[userId] != nil {
			// User STOPPED STREAMING but STILL IN VC
			delete(globals.StreamSessions, userId)
		} else if vs.ChannelID == "" && globals.StreamSessions[userId] == nil {
			// User left a voice channel
			if joinTime, ok := globals.VoiceSessions[userId]; ok {
				duration := time.Since(joinTime)
				err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(duration.Seconds()))
				if err != nil {
					fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
				}
				delete(globals.VoiceSessions, userId)
			}
		}
	}

}
