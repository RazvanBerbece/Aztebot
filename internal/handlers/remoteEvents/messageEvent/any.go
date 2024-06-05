package messageEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageCreatorUserId := m.Author.ID

	// Ignore all messages created by bots
	authorIsBot, err := member.IsBot(s, globalConfiguration.DiscordMainGuildId, messageCreatorUserId, false)
	if err != nil {
		return
	}
	if *authorIsBot {
		return
	}

	// Increase stats for user
	err = globalRepositories.UserStatsRepository.IncrementMessagesSentForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v\n", messageCreatorUserId, err)
	}

	err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v\n", messageCreatorUserId, err)
	}
	err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(messageCreatorUserId, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v\n", m.Author.ID, err)
	}

	// Publish experience grant message on the channel
	globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
		UserId: messageCreatorUserId,
		Points: globalConfiguration.ExperienceReward_MessageSent,
		Type:   "MSG_ACTIVITY",
	}

	// Award coins for activity
	go member.AwardFunds(s, messageCreatorUserId, 1*globalConfiguration.CoinReward_MessageSent, "MSG-SEND")

}
