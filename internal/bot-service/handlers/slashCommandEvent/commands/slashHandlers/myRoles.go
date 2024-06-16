package slashHandlers

import (
	"fmt"
	"log"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
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
				Content: "An error ocurred while trying to fetch your profile card.",
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

	usersRepository := repositories.NewUsersRepository()
	roles, err := usersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot display roles for user with id %s: %v", userId, err)
		return nil
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖    `%s`'s Roles", userName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	for _, role := range roles {
		var title string
		text := fmt.Sprintf("_%s_", role.Info)
		if role.Emoji != "" {
			// Role has an associated emoji
			title = fmt.Sprintf("`%s %s`", role.Emoji, role.DisplayName)
		} else {
			// Role doesn't have an associated emoji
			title = fmt.Sprintf("`%s`", role.DisplayName)
		}
		embed.
			AddField(title, text, false)
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
