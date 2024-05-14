package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
)

func HandleExperienceGrantEvents() {

	for xpEvent := range globalMessaging.ExperienceGrantsChannel {
		_, err := member.GrantMemberExperience(xpEvent.UserId, xpEvent.Points)
		if err != nil {
			fmt.Println("An error ocurred in the XP grant message handler for message", xpEvent, ":", err)
		}
	}

}
