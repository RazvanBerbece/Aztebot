package serverSlashHandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func HandleSlashAztebotHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userId := i.Interaction.Member.User.ID

	msg := sendHelpGuideToUser(userId)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func sendHelpGuideToUser(userId string) string {

	embedToSend := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("🤖   Command Guide").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Build guide message from all available and registered commands
	for _, cmd := range globalState.AztebotSlashCommands {
		title := fmt.Sprintf("`/%s`", cmd.Name)
		if len(cmd.Options) > 0 {
			for _, param := range cmd.Options {
				var required string
				if param.Required {
					required = "required"
				} else {
					required = "optional"
				}
				title += fmt.Sprintf(" `[%s (%s)]`", param.Name, required)
			}
		}

		if utils.StringInSlice(cmd.Name, globalConfiguration.RestrictedCommands) {
			// don't show restricted commands
			if member.IsStaff(userId, globalConfiguration.StaffRoles) {
				// unless a member of staff ran the /help handler
				enrichedTitle := fmt.Sprintf("%s *(higher staff command)*", title)
				embedToSend.AddField(enrichedTitle, cmd.Description, false)
			} else {
				continue
			}
		} else if utils.StringInSlice(cmd.Name, globalConfiguration.StaffCommands) {
			if member.IsStaff(userId, globalConfiguration.StaffRoles) {
				enrichedTitle := fmt.Sprintf("%s *(staff command)*", title)
				embedToSend.AddField(enrichedTitle, cmd.Description, false)
			} else {
				continue
			}
		} else {
			embedToSend.AddField(title, cmd.Description, false)
		}
	}

	paginationRow := embed.GetPaginationActionRowForEmbed(globalMessaging.PreviousPageOnEmbedEventId, globalMessaging.NextPageOnEmbedEventId)
	globalMessaging.DirectMessagesChannel <- events.DirectMessageEvent{
		UserId:        userId,
		Embed:         embedToSend,
		PaginationRow: &paginationRow,
	}

	return "You should have received a help guide for the `AzteBot` in your DMs."

}
