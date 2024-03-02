package profileSlashHandlers

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashYouRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	user, err := globalsRepo.UsersRepository.GetUser(targetUserId)
	if err != nil {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "An error ocurred while trying to fetch a user from the database")
		return
	}

	embed, err := RoleDisplayEmbedForUser(user.DiscordTag, targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while trying to fetch and display a user's roles: %v", err)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if embed == nil {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "An error ocurred while trying to fetch and display a user's roles.")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})
}
