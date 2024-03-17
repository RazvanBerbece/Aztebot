package timeoutSlashHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeouts(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeouts` command..."),
		},
	})

	activeTimeout, archived, err := member.GetMemberTimeouts(targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving timeouts for member with UID %s: `%s`", targetUserId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	user, err := s.User(targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Timeouts", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	if activeTimeout != nil {
		timeoutCreatedAtString := utils.FormatUnixAsString(activeTimeout.CreationTimestamp, "Mon, 02 Jan 2006 15:04:05 MST")

		// Format timeout expiry time
		timeoutCreationTime := time.Unix(activeTimeout.CreationTimestamp, 0)
		duration := time.Second * time.Duration(activeTimeout.SDuration)
		expiryTime := timeoutCreationTime.Add(duration)
		timeoutExpiryTimestampString := expiryTime.Format("Mon, 02 Jan 2006 15:04:05 MST")

		timeoutDesc := fmt.Sprintf("Given at `%s` \nfor reason `%s` \n(expiration date `%s`)", timeoutCreatedAtString, activeTimeout.Reason, timeoutExpiryTimestampString)
		embed.AddField("Active Timeout", timeoutDesc, false)
	}

	if len(archived) > 0 {
		archivedFieldName := "Archived Timeouts"
		archivedFieldValue := ""
		for idx, archivedTimeout := range archived {

			timeoutExpiryTime := time.Unix(archivedTimeout.ExpiryTimestamp, 0)
			timeoutExpiryTimeString := timeoutExpiryTime.Format("Mon, 02 Jan 2006 15:04:05 MST")

			archivedFieldValue += fmt.Sprintf("%d. `%s` | Expired: `%s` \n", idx+1, archivedTimeout.Reason, timeoutExpiryTimeString)

		}
		embed.AddField(archivedFieldName, archivedFieldValue, false)
	}

	if len(archived) <= 0 && activeTimeout == nil {
		embed.AddField("This member has no active or archived timeouts.", "", false)
	}

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
