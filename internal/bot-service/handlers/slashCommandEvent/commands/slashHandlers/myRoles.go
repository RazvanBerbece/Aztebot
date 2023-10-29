package slashHandlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMyRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: roleDisplayNameListForUser(i.Message.Member.User.ID),
		},
	})
}

func roleDisplayNameListForUser(userId string) string {
	return ""
}
