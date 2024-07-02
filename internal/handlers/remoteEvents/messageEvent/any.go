package messageEvent

import (
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
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
	var activityType = "MSG"
	globalMessaging.ActivityRegistrationsChannel <- events.ActivityEvent{
		UserId: messageCreatorUserId,
		Type:   &activityType,
	}

	// Publish experience grant message on the channel
	globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
		UserId: messageCreatorUserId,
		Points: globalConfiguration.ExperienceReward_MessageSent,
		Type:   "MSG_ACTIVITY",
	}

	// Award coins for activity
	globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
		UserId:   messageCreatorUserId,
		Funds:    1 * globalConfiguration.CoinReward_MessageSent,
		Activity: "MSG-SEND",
	}

}
