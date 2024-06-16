package profile

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashYou(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Gathering `/you` data..."),
		},
	})

	embed := GetProfileEmbedForUser(s, targetUserId)
	if embed == nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "An error ocurred while trying to fetch user's profile card.")
	}

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &embed,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)
}
