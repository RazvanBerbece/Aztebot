package channelHandlers

import (
	"fmt"
	"time"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandleMemberMessageDeletionEvents(s *discordgo.Session) {

	searchDepth := 5 // how many messages (n * 100) in each channel to check for deletion
	timeLimit := time.Hour * 24

	for deletionEvent := range globalMessaging.MessageDeletionChannel {
		err := member.DeleteMostRecentMemberMessages(s, deletionEvent.GuildId, deletionEvent.UserId, searchDepth, timeLimit)
		if err != nil {
			fmt.Println("An error ocurred in the member deletion event handler for message", deletionEvent, ":", err)
		}
	}

}
