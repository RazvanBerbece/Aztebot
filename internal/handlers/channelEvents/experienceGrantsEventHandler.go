package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandleExperienceGrantEvents(s *discordgo.Session, logger logging.Logger) {

	for xpEvent := range globalMessaging.ExperienceGrantsChannel {
		_, err := member.GrantMemberExperience(xpEvent.UserId, xpEvent.Points)
		if err != nil {
			fmt.Println("An error ocurred in the XP grant message handler for message", xpEvent, ":", err)

			logMsg := fmt.Sprintf("Failed to grant `%f` XP to `%s`", xpEvent.Points, xpEvent.UserId)
			go logger.LogError(logMsg)
		}
	}

}
