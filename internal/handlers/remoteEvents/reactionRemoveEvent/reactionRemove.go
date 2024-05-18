package messageEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func ReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {

	// Ignore all reactions removed by the bot itself
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

	err = globalRepositories.UserStatsRepository.DecrementReactionsReceivedForUser(messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) reaction count: %v\n", messageOwnerUid, err)
	}

	// Remove experience points from message owner
	currentXp, err := member.RemoveMemberExperience(messageOwnerUid, "REACT_REWARD")
	if err != nil {
		fmt.Printf("An error ocurred while removing reaction received rewards (%d) from user (%s): %v\n", currentXp, messageOwnerUid, err)
	}

}
