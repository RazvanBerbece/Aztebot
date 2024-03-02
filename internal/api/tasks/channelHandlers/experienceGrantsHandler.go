package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
)

func HandleExperienceGrantsMessages(debug bool) {

	for msg := range globals.ExperienceGrantsChannel {
		_, err := member.GrantMemberExperience(msg.UserId, msg.Activity, msg.Points)
		if err != nil {
			fmt.Println("An error ocurred in the XP grant message handler for message", msg, ":", err)
		}
		if debug {
			fmt.Println(msg)
		}
	}

}
