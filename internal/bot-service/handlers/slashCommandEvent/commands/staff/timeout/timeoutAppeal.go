package timeoutSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeoutAppeal(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userId := i.Member.User.ID

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeout_appeal` command..."),
		},
	})

	user, err := s.User(userId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID %s: %v", userId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
	}

	// Timeout appeal logic
	err = member.AppealTimeout(s, globals.DiscordMainGuildId, userId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while appealing timeout for user with UID `%s`: %v", userId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Timeout Appeal from `%s`", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
