package jailSlashHandlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleSlashUnjail(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/unjail` command..."),
	// 	},
	// })

	// jailedUser, user, err := member.UnjailMember(s, globals.DiscordMainGuildId, targetUserId)
	// if err != nil {
	// 	utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
	// 	return
	// }

	// embed := embed.NewEmbed().
	// 	SetTitle(fmt.Sprintf("🤖⚠️   Unjailed `%s` !", user.DiscordTag)).
	// 	SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
	// 	SetColor(000000).
	// 	AddField("Reason", jailedUser.Reason, false).
	// 	AddField("Given Task", jailedUser.Reason, false).
	// 	AddField("Timestamp", utils.FormatUnixAsString(jailedUser.JailedAt, "Mon, 02 Jan 2006 15:04:05 MST"), false)

	// Final response
	// editContent := ""
	// editWebhook := discordgo.WebhookEdit{
	// 	Content: &editContent,
	// 	Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	// }
	// s.InteractionResponseEdit(i.Interaction, &editWebhook)

	return

}
