package channelHandlers

import (
	"database/sql"
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
			if err == sql.ErrNoRows {
				// Ignore adding XP in cases where the member left the server in the meantime
				continue
			}
			fmt.Println("An error ocurred in the XP grant message handler for message", xpEvent, ":", err)
			logMsg := fmt.Sprintf("Failed to grant `%f` XP (`%s`) to `%s`\nTrace: %s", xpEvent.Points, xpEvent.Type, xpEvent.UserId, err)
			go logger.LogError(logMsg)
		}
	}

}
