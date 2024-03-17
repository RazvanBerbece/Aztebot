package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
)

func HandleExperienceGrantEvents() {

	for xpEvent := range globalMessaging.ExperienceGrantsChannel {
		_, err := member.GrantMemberExperience(xpEvent.UserId, xpEvent.Activity, xpEvent.Points)
		if err != nil {
			fmt.Println("An error ocurred in the XP grant message handler for message", xpEvent, ":", err)
		}
	}

}
