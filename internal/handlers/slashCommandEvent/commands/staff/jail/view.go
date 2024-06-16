package jailSlashHandlers

import (
	"fmt"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashJailView(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/jail-view` command..."),
		},
	})

	jailed, err := globalRepositories.JailRepository.GetJail()
	if err != nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
		return
	}

	// Build response embed with a detailed jail view
	embed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️   The OTA Jail").
		SetDescription(fmt.Sprintf("The OTA Jail is the place where the convicted server members are sent to.\nCurrently, there are `%d` members imprisoned in the OTA Jail.", len(jailed))).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddLineBreakField().
		AddField("Imprisoned Members List", "", false)

	for idx, jailedUser := range jailed {
		user, err := globalRepositories.UsersRepository.GetUser(jailedUser.UserId)
		if err != nil {
			utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
			return
		}
		embed.AddField("", fmt.Sprintf("%d. `%s`\nConvicted on: `%s`, for the following reason: `%s`\nHas to complete this task for release: `%s`", idx+1, user.DiscordTag, utils.FormatUnixAsString(jailedUser.JailedAt, "Mon, 02 Jan 2006 15:04:05 MST"), jailedUser.Reason, jailedUser.TaskToComplete), false)
	}

	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
