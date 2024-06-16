package messageEvent

import (
	"fmt"
	"time"

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

}
