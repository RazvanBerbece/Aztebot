package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
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
	authorIsBot, err := member.MemberIsBot(s, globals.DiscordMainGuildId, messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while checking against bot application: %v\n", err)
	}
	if authorIsBot == nil {
		return
	}
	if *authorIsBot {
		return
	}

	err = globalsRepo.UserStatsRepository.DecrementReactionsReceivedForUser(messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) reaction count: %v", messageOwnerUid, err)
	}

	// Remove experience points from message owner
	currentXp, err := member.RemoveMemberExperience(messageOwnerUid, "REACT_REWARD")
	if err != nil {
		fmt.Printf("An error ocurred while removing reaction received rewards (%d) from user (%s): %v", currentXp, messageOwnerUid, err)
	}

}
