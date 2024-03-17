package profileSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashSetGender(s *discordgo.Session, i *discordgo.InteractionCreate) {

	genderInput := i.ApplicationCommandData().Options[0].StringValue()

	// Dirty Hack 25 Feb 2024
	// It seems that it's not straightforward at all to get the display name of the argument option,
	// so we resort to this for the meantime to get a nicely looking activity and multiplier name
	genderName := getArgumentDisplayNames(genderInput)

	commandOwnerUserId := i.Member.User.ID
	commandOwnerUsername := i.Member.User.Username

	err := member.SetGender(commandOwnerUserId, genderInput)
	if err != nil {
		fmt.Printf("An error ocurred while setting gender for user with UID %s: %v", commandOwnerUserId, err)
		utils.SendErrorEmbedResponse(s, i.Interaction, fmt.Sprintf("Failed to set gender for `%s`", commandOwnerUsername))
		return
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   Updated Profile Gender For `%s`", commandOwnerUsername)).
		SetColor(000000).
		AddField(fmt.Sprintf("Updated to `%s`.", genderName), "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func getArgumentDisplayNames(genderInput string) string {

	var genderName string

	switch genderInput {
	case "male":
		genderName = "Male"
	case "female":
		genderName = "Female"
	case "nonbin":
		genderName = "Nonbinary"
	case "other":
		genderName = "Other"
	}

	return genderName
}
