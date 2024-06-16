package messageEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
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
	authorIsBot, err := member.MemberIsBot(s, globals.DiscordMainGuildId, messageOwnerUid)
	if err != nil {
		return
	}
	if authorIsBot == nil {
		return
	}
	if *authorIsBot {
		return
	}

	err = globalsRepo.UserStatsRepository.IncrementReactionsReceivedForUser(messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) reaction count: %v", messageOwnerUid, err)
	}

	err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(r.UserID)
	if err != nil {
		fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", r.UserID, err)
	}
	err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(r.UserID, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", r.UserID, err)
	}

	// Grant experience points
	currentXp, err := member.GrantMemberExperience(messageOwnerUid, "REACT_REWARD", nil)
	if err != nil {
		fmt.Printf("An error ocurred while granting reaction received rewards (%d) to user (%s): %v", currentXp, messageOwnerUid, err)
	}

}
