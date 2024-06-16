package massPingSlashHandlers

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMassDm(s *discordgo.Session, i *discordgo.InteractionCreate) {

	commandOwnerUserId := i.Member.User.ID

	// This is a higher staff command however restrict specifically anyone lower than "Developer" from using it
	ownerStaffRole, err := member.GetMemberStaffRole(commandOwnerUserId, globalConfiguration.StaffRoles)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't check for command owner permissions: %v", err)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if ownerStaffRole == nil || (ownerStaffRole.DisplayName != "Developer" && ownerStaffRole.DisplayName != "Dominus" && ownerStaffRole.DisplayName != "Arhitect") {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "Command owner doesn't have the right permissions to use this command.")
		return
	}

	msg := i.ApplicationCommandData().Options[0].StringValue()

	// Send response embed
	embed := embed.NewEmbed().
		SetColor(000000).
		SetAuthor("AzteBot Mass DM", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetDescription(msg)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}
