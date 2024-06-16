package voiceStateUpdateEvent

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
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
	if newChannelName, isCreateChannelCommand := globalConfiguration.DynamicChannelCreateButtonIds[vs.ChannelID]; isCreateChannelCommand {

		// Privacy status of channel (public or private)
		channelIsPrivate := false
		if strings.Contains(newChannelName, "Private") {
			channelIsPrivate = true
		}

		// Publish channel creation event
		globalMessaging.ChannelCreationsChannel <- events.VoiceChannelCreateEvent{
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
		if _, isAfkChannel := globalConfiguration.AfkChannels[vs.ChannelID]; isAfkChannel {
			// Don't register audio sessions in the AFK zones
			delete(globalState.MusicSessions, userId)
			delete(globalState.VoiceSessions, userId)
			delete(globalState.StreamSessions, userId)
			delete(globalState.DeafSessions, userId)
			return
		}
	}

	if vs.SelfStream && vs.ChannelID != "" {
		// User STARTED STREAMING
		now := time.Now()
		globalState.StreamSessions[userId] = &now

		err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
		}

		err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(userId, now.Unix())
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", userId, err)
		}
	} else if vs.SelfStream && vs.ChannelID == "" {
		// The Discord API does something weird and sends SelfStream as true
		// when a user leaves a VC directly without stopping streaming first
		if joinTime, ok := globalState.VoiceSessions[userId]; ok {

			duration := time.Since(joinTime)
			secondsSpent := duration.Seconds()

			err := globalRepositories.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(secondsSpent))
			if err != nil {
				fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
			}

			globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
				UserId: userId,
				Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
				Type:   "VOICE_ACTIVITY",
			}

			globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
				UserId:   userId,
				Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
				Activity: "TIME-VC",
			}

			delete(globalState.VoiceSessions, userId)
			delete(globalState.StreamSessions, userId)
			delete(globalState.MusicSessions, userId)
		}
	} else {
		if vs.ChannelID != "" && globalState.StreamSessions[userId] == nil && globalState.MusicSessions[userId] == nil {
			// User JOINED a VC but NOT STREAMING
			if utils.TargetChannelIsForMusicListening(globalConfiguration.MusicChannels, vs.ChannelID) {
				now := time.Now()
				globalState.MusicSessions[userId] = map[string]*time.Time{
					vs.ChannelID: &now,
				}
			} else {
				globalState.VoiceSessions[userId] = time.Now()
			}

			err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(userId)
			if err != nil {
				fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", userId, err)
			}

			err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(userId, time.Now().Unix())
			if err != nil {
				fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", userId, err)
			}
		} else if vs.ChannelID != "" && globalState.StreamSessions[userId] != nil {
			delete(globalState.StreamSessions, userId)
		} else if vs.ChannelID == "" && globalState.StreamSessions[userId] == nil {
			// User LEFT THE VOICE CHANNEL
			musicSession, userHadMusicSession := globalState.MusicSessions[userId]
			if userHadMusicSession {
				// User was on a music channel
				for _, joinTime := range musicSession {

					duration := time.Since(*joinTime)
					secondsSpent := duration.Seconds()

					err := globalRepositories.UserStatsRepository.AddToTimeSpentListeningMusic(userId, int(secondsSpent))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", userId, err)
					}

					globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
						UserId: userId,
						Points: globalConfiguration.ExperienceReward_InMusic * secondsSpent,
						Type:   "MUSIC_ACTIVITY",
					}

					globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
						UserId:   userId,
						Funds:    globalConfiguration.CoinReward_InMusic * secondsSpent,
						Activity: "TIME-MUSIC",
					}
				}
			} else {
				// User was on any other VC
				if joinTime, ok := globalState.VoiceSessions[userId]; ok {

					duration := time.Since(joinTime)
					secondsSpent := duration.Seconds()

					err := globalRepositories.UserStatsRepository.AddToTimeSpentInVoiceChannels(userId, int(duration.Seconds()))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", userId, err)
					}

					globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
						UserId: userId,
						Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
						Type:   "VOICE_ACTIVITY",
					}

					globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
						UserId:   userId,
						Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
						Activity: "TIME-VC",
					}
				}
			}
			delete(globalState.MusicSessions, userId)
			delete(globalState.VoiceSessions, userId)
			delete(globalState.StreamSessions, userId)
			delete(globalState.DeafSessions, userId)
		}
	}

}
