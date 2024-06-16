package warning

import (
	"fmt"
	"time"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashWarns(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Gathering your `/warns` data..."),
		},
	})

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving user to get warnings for: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId),
			},
		})
	}

	// Retrieve all warnings for user with given UID
	warns, err := globalsRepo.WarnsRepository.GetWarningsForUser(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving user warns: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error ocurred while retrieving user warnings for user with ID %s provided in the slash command.", targetUserId),
			},
		})
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   `%s`'s Warnings", user.Username)).
		SetDescription(fmt.Sprintf("`%s` has %d warnings.", user.Username, len(warns))).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	for idx, warn := range warns {
		// Format CreatedAt
		var warnCreatedAt time.Time
		var warnCreatedAtString string
		warnCreatedAt = time.Unix(warn.CreationTimestamp, 0).UTC()
		warnCreatedAtString = warnCreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")

		warningTitle := fmt.Sprintf("Warning `#%d`", idx+1)
		warningReason := fmt.Sprintf("Reason: `%s`", warn.Reason) + "    *(created at " + warnCreatedAtString + ")*"
		embed.AddField(warningTitle, warningReason, false)
	}

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
