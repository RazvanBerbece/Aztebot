package member

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func SendDirectMessageToMember(s *discordgo.Session, userId string, msg string) error {
	errDm := dm.DmUser(s, userId, msg)
	if errDm != nil {
		fmt.Printf("Error sending DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}

func SendDirectEmbedToMember(s *discordgo.Session, userId string, embed embed.Embed) error {
	errDm := dm.DmEmbedUser(s, userId, *embed.MessageEmbed)
	if errDm != nil {
		fmt.Printf("Error sending embed DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}
