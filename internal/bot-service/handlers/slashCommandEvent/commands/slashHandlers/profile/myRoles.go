package profile

import (
	"fmt"
	"log"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMyRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Attempt a sync
	err := ProcessUserUpdate(i.Interaction.Member.User.ID, s, i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your roles.",
			},
		})
	}

	embed := roleDisplayEmbedForUser(i.Interaction.Member.User.Username, i.Interaction.Member.User.ID)
	if embed == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your roles.",
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: roleDisplayEmbedForUser(i.Interaction.Member.User.Username, i.Interaction.Member.User.ID),
		},
	})
}

func roleDisplayEmbedForUser(userName string, userId string) []*discordgo.MessageEmbed {

	roles, err := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot display roles for user with id %s: %v", userId, err)
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Roles", userName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
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

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
