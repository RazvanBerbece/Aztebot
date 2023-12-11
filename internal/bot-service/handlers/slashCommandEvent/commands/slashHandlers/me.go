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

	embed := displayEmbedForUser(s, i.Interaction.Member.User.ID)
	if embed == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your profile card.",
			},
		})
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})
}

func displayEmbedForUser(s *discordgo.Session, userId string) []*discordgo.MessageEmbed {

	usersRepository := repositories.NewUsersRepository()
	user, err := usersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
		return nil
	}

	// Format CreatedAt
	userCreatedTime := time.Unix(*user.CreatedAt, 0).UTC()
	userCreatedTimeString := userCreatedTime.Format("January 2, 2006")

	// Process highest role
	var highestRole dataModels.Role
	roles, err := usersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		return nil
	}
	highestRole = roles[len(roles)-1] // role IDs for users are stored in DB in ascending order by rank, so the last one is the highest

	// Get the profile picture url
	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		return nil
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖    `%s`'s Profile Card", user.DiscordTag)).
		SetDescription(fmt.Sprintf("`%s CIRCLE`", user.CurrentCircle)).
		SetThumbnail(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userId, apiUser.Avatar)).
		SetColor(000000).
		AddField(fmt.Sprintf("Aztec since:  `%s`", userCreatedTimeString), "", false).
		AddField(fmt.Sprintf("Highest obtained role:  `%s`", highestRole.DisplayName), "", false)

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
