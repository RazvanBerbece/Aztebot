package slashHandlers

import (
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: displayEmbedForUser(i.Interaction.Member.User.ID),
		},
	})
}

func displayEmbedForUser(userId string) []*discordgo.MessageEmbed {

	usersRepository := repositories.NewUsersRepository()
	user, err := usersRepository.GetUser(userId)
	if err != nil {
		log.Fatalf("Cannot retrieve user with id %s: %v", userId, err)
	}

	// Format CreatedAt
	userCreatedTime := time.Unix(*user.CreatedAt, 0).UTC()
	userCreatedTimeString := userCreatedTime.Format("January 2, 2006")

	// Process highest role
	var highestRole dataModels.Role
	roles, err := usersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Fatalf("Cannot retrieve roles for user with id %s: %v", userId, err)
	}
	highestRole = roles[len(roles)-1] // role IDs for users are stored in DB in ascending order by rank, so the last one is the highest

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖    `%s`'s Profile Card", user.DiscordTag)).
		SetDescription(fmt.Sprintf("`%s CIRCLE`", user.CurrentCircle)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		// Add field for highest role obtained
		// Add field for date created / verified
		AddField(fmt.Sprintf("Aztec since:  `%s`", userCreatedTimeString), "", false).
		AddField(fmt.Sprintf("Highest obtained role:  `%s %s`", highestRole.Emoji, highestRole.DisplayName), "", false)

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
