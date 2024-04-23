package profileSlashHandlers

import (
	"fmt"
	"log"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMyRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {

	embed, err := RoleDisplayEmbedForUser(i.Interaction.Member.User.Username, i.Interaction.Member.User.ID)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while trying to fetch and display your roles: %v", err)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if embed == nil {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "An error ocurred while trying to fetch and display your roles.")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})
}

func RoleDisplayEmbedForUser(userName string, userId string) ([]*discordgo.MessageEmbed, error) {

	roles, err := globalRepositories.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot display roles for user with id %s: %v", userId, err)
		return nil, err
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Roles", userName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000)

	for _, role := range roles {
		var title string
		var text string = role.Info
		if role.Emoji != "" {
			// Role has an associated emoji
			title = fmt.Sprintf("`%s`", role.DisplayName)
		} else {
			// Role doesn't have an associated emoji
			title = fmt.Sprintf("`%s`", role.DisplayName)
		}
		// Only add field for role description if there is a description available
		if text != "" {
			text = fmt.Sprintf("_%s_", text) // italic
			embed.
				AddField(title, text, false)
		} else {
			embed.
				AddField(title, "", false)
		}
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}, nil
}
