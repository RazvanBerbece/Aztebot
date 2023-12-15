package messageEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
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

	userStatsRepo := repositories.NewUsersStatsRepository()
	err = userStatsRepo.IncrementReactionsReceivedForUser(messageOwnerUid)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) reaction count: %v", messageOwnerUid, err)
	}

}
