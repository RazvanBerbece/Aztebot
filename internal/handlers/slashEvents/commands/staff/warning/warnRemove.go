package warningSlashHandlers

import (
	"database/sql"
	"fmt"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashWarnRemoveOldest(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	warn, err := RemoveWarningFromUser(s, i, targetUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("User with ID %s has no warnings to remove.", targetUserId),
				},
			})
			return
		} else {
			fmt.Printf("An error ocurred while removing warning from user: %v\n", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("An error ocurred while removing warning from user with ID %s.", targetUserId),
				},
			})
			return
		}
	}

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while sending warning removal embed response: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while sending warning removal embed response.",
			},
		})
		return
	}

	// Format CreatedAt
	warnCreatedAtString := utils.FormatUnixAsString(warn.CreationTimestamp, "Mon, 02 Jan 2006 15:04:05 MST")

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖🔨   Warning removed from `%s`", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Reason", warn.Reason, false).
		AddField("Timestamp", warnCreatedAtString, false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func RemoveWarningFromUser(s *discordgo.Session, i *discordgo.InteractionCreate, userId string) (*dax.Warn, error) {

	warn, err := globalRepositories.WarnsRepository.GetOldestWarnForUser(userId)
	if err != nil {
		fmt.Printf("Error occured while getting oldest warning for user %s: %v\n", userId, err)
		return nil, err
	}

	err = globalRepositories.WarnsRepository.DeleteOldestWarningForUser(userId)
	if err != nil {
		fmt.Printf("Error occured while deleting oldest warnings for user: %v\n", err)
		return nil, err
	}

	return warn, nil
}
