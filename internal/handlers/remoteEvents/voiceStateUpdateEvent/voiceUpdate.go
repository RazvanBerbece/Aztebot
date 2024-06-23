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
			ParentMemberId:  vs.UserID,
			ParentGuildId:   member.GuildID,
		}

		return
	}

	// Don't register audio sessions in the AFK zones
	if vs.ChannelID != "" {
		if _, isAfkChannel := globalConfiguration.AfkChannels[vs.ChannelID]; isAfkChannel {
			delete(globalState.MusicSessions, vs.UserID)
			delete(globalState.VoiceSessions, vs.UserID)
			delete(globalState.StreamSessions, vs.UserID)
			return
		}
	}

	// Don't register audio sessions when users are deafened
	if vs.ChannelID != "" {
		if vs.SelfDeaf || vs.Deaf {
			// TODO: Store leftover durations
			// And remove user from active audio sessions (they'll be added back in when they un-deaf)
			delete(globalState.MusicSessions, vs.UserID)
			delete(globalState.VoiceSessions, vs.UserID)
			delete(globalState.StreamSessions, vs.UserID)
			return
		}
	}

	// Track durations for new stream sessions
	if vs.ChannelID != "" && vs.SelfStream {

		now := time.Now()
		globalState.StreamSessions[vs.UserID] = &now

		err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(vs.UserID)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", vs.UserID, err)
		}

		err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(vs.UserID, now.Unix())
		if err != nil {
			fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", vs.UserID, err)
		}
	} else if vs.ChannelID == "" && vs.SelfStream {
		// The Discord API does something weird and sends SelfStream as true
		// when a user leaves a VC directly without stopping streaming first
		if joinTime, ok := globalState.VoiceSessions[vs.UserID]; ok {

			duration := time.Since(joinTime)
			secondsSpent := duration.Seconds()

			err := globalRepositories.UserStatsRepository.AddToTimeSpentInVoiceChannels(vs.UserID, int(secondsSpent))
			if err != nil {
				fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", vs.UserID, err)
			}

			globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
				UserId: vs.UserID,
				Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
				Type:   "VOICE_ACTIVITY",
			}

			globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
				UserId:   vs.UserID,
				Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
				Activity: "TIME-VC",
			}

			delete(globalState.VoiceSessions, vs.UserID)
			delete(globalState.StreamSessions, vs.UserID)
			delete(globalState.MusicSessions, vs.UserID)
		}
	} else {
		// Track duration in simple voice session
		if vs.ChannelID != "" && globalState.StreamSessions[vs.UserID] == nil && globalState.MusicSessions[vs.UserID] == nil {
			if utils.TargetChannelIsForMusicListening(globalConfiguration.MusicChannels, vs.ChannelID) {
				now := time.Now()
				globalState.MusicSessions[vs.UserID] = map[string]*time.Time{
					vs.ChannelID: &now,
				}
			} else {
				globalState.VoiceSessions[vs.UserID] = time.Now()
			}

			err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(vs.UserID)
			if err != nil {
				fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", vs.UserID, err)
			}

			err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(vs.UserID, time.Now().Unix())
			if err != nil {
				fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v", vs.UserID, err)
			}
		} else if vs.ChannelID != "" && globalState.StreamSessions[vs.UserID] != nil {
			delete(globalState.StreamSessions, vs.UserID)
		} else if vs.ChannelID == "" && globalState.StreamSessions[vs.UserID] == nil { // Track duration for voice session or music leaver

			musicSession, userHadMusicSession := globalState.MusicSessions[vs.UserID]

			if userHadMusicSession {
				for _, joinTime := range musicSession {

					duration := time.Since(*joinTime)
					secondsSpent := duration.Seconds()

					err := globalRepositories.UserStatsRepository.AddToTimeSpentListeningMusic(vs.UserID, int(secondsSpent))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent listening music for user with id %s: %v", vs.UserID, err)
					}

					globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
						UserId: vs.UserID,
						Points: globalConfiguration.ExperienceReward_InMusic * secondsSpent,
						Type:   "MUSIC_ACTIVITY",
					}

					globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
						UserId:   vs.UserID,
						Funds:    globalConfiguration.CoinReward_InMusic * secondsSpent,
						Activity: "TIME-MUSIC",
					}
				}
			} else { // User was on any other VC
				if joinTime, ok := globalState.VoiceSessions[vs.UserID]; ok {

					duration := time.Since(joinTime)
					secondsSpent := duration.Seconds()

					err := globalRepositories.UserStatsRepository.AddToTimeSpentInVoiceChannels(vs.UserID, int(duration.Seconds()))
					if err != nil {
						fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", vs.UserID, err)
					}

					globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
						UserId: vs.UserID,
						Points: globalConfiguration.ExperienceReward_InVc * secondsSpent,
						Type:   "VOICE_ACTIVITY",
					}

					globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
						UserId:   vs.UserID,
						Funds:    globalConfiguration.CoinReward_InVc * secondsSpent,
						Activity: "TIME-VC",
					}
				}
			}

			delete(globalState.MusicSessions, vs.UserID)
			delete(globalState.VoiceSessions, vs.UserID)
			delete(globalState.StreamSessions, vs.UserID)
		}
	}

}
