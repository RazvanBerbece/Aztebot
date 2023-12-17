package slashHandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
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

	// Build guide message from all available and registered commands
	helpGuideMsg := "Current list of all available commands for the AzteBot: \n\n"
	for _, cmd := range globals.AztebotRegisteredCommands {
		helpGuideMsg += "- " + fmt.Sprintf("`%s`", cmd.Name) + " -> " + cmd.Description + "\n\n"
	}

	errDm := dm.SendHelpDmToUser(s, i, userId, helpGuideMsg)
	if errDm != nil {
		fmt.Println("Error sending DM: ", errDm)
		return "An error occured while DMing you the help guide."
	}
	return "You should have received a help guide for the `@AzteBot` in your DMs."

}
