package timeoutSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeoutRemoveActive(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeout_remove_active` command..."),
		},
	})

	user, err := s.User(targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	activeTimeout, _, err := member.GetMemberTimeouts(targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving timeout data for user with ID %s: %v", targetUserId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}
	if activeTimeout == nil {
		errMsg := fmt.Sprintf("No active timeout was found for user with ID `%s`", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	err = member.ClearMemberActiveTimeout(s, globals.DiscordMainGuildId, targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while clearing timeout for user with ID %s: %v", targetUserId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Timeout cleared from `%s`", user.Username)).
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
