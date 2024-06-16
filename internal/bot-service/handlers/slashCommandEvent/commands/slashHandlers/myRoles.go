package slashHandlers

import (
	"log"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMyRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: roleDisplayNameListForUser(i.Interaction.Member.User.ID),
		},
	})
}

func roleDisplayNameListForUser(userId string) string {
	usersRepository := repositories.NewUsersRepository()
	roles, err := usersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Fatalf("Cannot display role names for user with id %s: %v", userId, err)
	}

	displayString := "Your roles are: "
	for index, role := range roles {
		if index == len(roles)-1 {
			displayString += role.DisplayName
			break
		}
		displayString += role.DisplayName + ", "
	}

	return displayString
}
