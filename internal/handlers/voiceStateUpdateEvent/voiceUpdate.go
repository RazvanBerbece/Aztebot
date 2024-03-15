package voiceStateUpdateEvent

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
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

	if member.User.Bot {
		// Ignore bot voice state updates
		return
	}

	userId := member.User.ID

	// If a dynamic room creation command
	if newChannelName, isCreateChannelCommand := globals.DynamicChannelCreateButtonIds[vs.ChannelID]; isCreateChannelCommand {

		// Privacy status of channel (public or private)
		channelIsPrivate := false
		if strings.Contains(newChannelName, "Private") {
			channelIsPrivate = true
		}

		// Publish channel creation event
		globals.ChannelCreationsChannel <- events.VoiceChannelCreateEvent{
			Name:            newChannelName,
			Private:         channelIsPrivate,
			ParentChannelId: vs.ChannelID,
			Description:     "This is a dynamically generated voice channel!",
			ParentMemberId:  userId,
			ParentGuildId:   member.GuildID,
		}

		return
	}

	if vs.ChannelID != "" {
		if _, isAfkChannel := globals.AfkChannels[vs.ChannelID]; isAfkChannel {
			// Don't register audio sessions in the AFK zones
			delete(globals.MusicSessions, userId)
			delete(globals.VoiceSessions, userId)
			delete(globals.StreamSessions, userId)
			delete(globals.DeafSessions, userId)
			return
		}
	}

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

			globals.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
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
			if utils.TargetChannelIsForMusicListening(globals.MusicChannels, vs.ChannelID) {
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

					globals.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
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

					globals.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
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
