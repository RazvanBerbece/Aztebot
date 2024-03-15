package jailSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashUnjail(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	if !utils.IsValidDiscordUserId(targetUserId) {
		errMsg := fmt.Sprintf("The provided `user` command argument is invalid. (term: `%s`)", targetUserId)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/unjail` command..."),
		},
	})

	var err error
	var user *dataModels.User
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-jail"]; channelExists {
		_, user, err = member.UnjailMember(s, globalConfiguration.DiscordMainGuildId, targetUserId, globalConfiguration.JailedRoleName, channel.ChannelId)
		if err != nil {
			utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
			return
		}
	} else {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "The `/unjail` command cannot be used without a designated notification channel to send various jail related alerts to.")
		return
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Unjailed `%s` !", user.DiscordTag)).
		SetColor(000000)

	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
