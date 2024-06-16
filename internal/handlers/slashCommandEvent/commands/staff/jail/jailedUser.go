package jailSlashHandlers

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashJailedUser(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	// Input validation
	if !utils.IsValidDiscordUserId(targetUserId) {
		errMsg := fmt.Sprintf("The provided `user` command argument is invalid. (term: `%s`)", targetUserId)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/jailed-user` command..."),
		},
	})

	userIsInJail := globalsRepo.JailRepository.UserIsJailed(targetUserId)
	if userIsInJail <= 0 {
		if userIsInJail == -1 {
			utils.ErrorEmbedResponseEdit(s, i.Interaction, "an error ocurred while checking if the given user is currently in jail")
			return
		}
		utils.ErrorEmbedResponseEdit(s, i.Interaction, fmt.Sprintf("user `#%s` was not found in the OTA jail", targetUserId))
		return
	}

	jailedUser, err := globalsRepo.JailRepository.GetJailedUser(targetUserId)
	if err != nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
		return
	}

	// Build response embed with a detailed jail view
	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("👮🏽‍♀️⛓️   Jailed User Record - `#%s`", jailedUser.UserId)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Convicted since", fmt.Sprintf("`%s`", utils.FormatUnixAsString(jailedUser.JailedAt, "Mon, 02 Jan 2006 15:04:05 MST")), false).
		AddField("Charges", fmt.Sprintf("`%s`", jailedUser.Reason), false).
		AddField("Assigned task", fmt.Sprintf("`%s`", jailedUser.TaskToComplete), false)

	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
