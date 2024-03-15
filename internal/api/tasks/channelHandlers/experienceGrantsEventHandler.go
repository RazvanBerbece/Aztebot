package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
)

func HandleExperienceGrantEvents() {

	for xpEvent := range globals.ExperienceGrantsChannel {
		_, err := member.GrantMemberExperience(xpEvent.UserId, xpEvent.Activity, xpEvent.Points)
		if err != nil {
			fmt.Println("An error ocurred in the XP grant message handler for message", xpEvent, ":", err)
		}
	}

}
