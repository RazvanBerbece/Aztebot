package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	// Ignore all reactions created by the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Println("Error retrieving message:", err)
		return
	}

	messageOwnerUid := message.Author.ID

	// Ignore all messages created by bots
	authorIsBot, err := member.IsBot(s, globalConfiguration.DiscordMainGuildId, messageOwnerUid, false)
	if err != nil {
		return
	}
	if authorIsBot == nil {
		return
	}
	if *authorIsBot {
		return
	}

	err = globalRepositories.UserStatsRepository.IncrementReactionsReceivedForUser(messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) reaction count: %v", messageOwnerUid, err)
	}

	globalMessaging.ActivityRegistrationsChannel <- events.ActivityEvent{
		UserId: r.UserID,
	}

	globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
		UserId: messageOwnerUid,
		Points: globalConfiguration.ExperienceReward_ReactionReceived,
		Type:   "REACT_ACTIVITY",
	}

	// Award coins for activity
	globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
		GuildId:  r.GuildID,
		UserId:   messageOwnerUid,
		Funds:    1 * globalConfiguration.CoinReward_ReactionReceived,
		Activity: "REACT-RECV",
	}

}
