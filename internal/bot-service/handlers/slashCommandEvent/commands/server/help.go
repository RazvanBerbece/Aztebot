package serverSlashHandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func HandleSlashAztebotHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userId := i.Interaction.Member.User.ID

	msg := sendHelpGuideToUser(s, i, userId)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func sendHelpGuideToUser(s *discordgo.Session, i *discordgo.InteractionCreate, userId string) string {

	embed := embed.NewEmbed().
		SetTitle("🤖   Command Guide").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Build guide message from all available and registered commands
	for _, cmd := range globals.AztebotRegisteredCommands {

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

		if utils.StringInSlice(cmd.Name, globals.RestrictedCommands) || utils.StringInSlice(cmd.Name, globals.StaffCommands) {
			// If a restricted or staff command, do not show
			if member.IsStaff(userId) {
				// unless a member of staff executed the command
				embed.AddField(fmt.Sprintf("%s *(staff command)*", title), cmd.Description, false)
			} else {
				continue
			}
		} else {
			embed.AddField(title, cmd.Description, false)
		}
	}

	errDm := dm.SendEmbedToUser(s, i, userId, embed.MessageEmbed)
	if errDm != nil {
		fmt.Println("Error sending DM: ", errDm)
		return "An error occured while DMing you the help guide."
	}
	return "You should have received a help guide for the `AzteBot` in your DMs."

}
