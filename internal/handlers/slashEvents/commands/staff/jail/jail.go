package jailSlashHandlers

import (
	"fmt"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashJail(s *discordgo.Session, i *discordgo.InteractionCreate) {

	commandOwnerUserId := i.Member.User.ID

	// This is a staff command however restrict specifically anyone lower than "Moderator" from using it
	ownerStaffRole, err := member.GetMemberStaffRole(commandOwnerUserId, globalConfiguration.StaffRoles)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't check for command owner permissions: %v", err)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if ownerStaffRole == nil || ownerStaffRole.DisplayName == "Trial Moderator" {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "Command owner doesn't have the right permissions to use this command.")
		return
	}

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	reason := i.ApplicationCommandData().Options[1].StringValue()

	if !utils.IsValidDiscordUserId(targetUserId) {
		errMsg := fmt.Sprintf("The provided `user` command argument is invalid. (term: `%s`)", targetUserId)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/jail` command..."),
		},
	})

	var user *dax.User
	var jailedUser *dax.JailedUser
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-jail"]; channelExists {
		jailedUser, user, err = member.JailMember(s, globalConfiguration.DiscordMainGuildId, targetUserId, reason, globalConfiguration.JailedRoleName, channel.ChannelId)
		if err != nil {
			utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
			return
		}
	} else {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "The `/jail` command cannot be used without a designated notification channel to send various jail related alerts to.")
		return
	}

	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖⚠️   Jailed `%s`", user.DiscordTag)).
		SetColor(000000).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		AddField("Reason", jailedUser.Reason, false).
		AddField("Given Task", jailedUser.TaskToComplete, false).
		AddField("Timestamp", utils.FormatUnixAsString(jailedUser.JailedAt, "Mon, 02 Jan 2006 15:04:05 MST"), false)

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
