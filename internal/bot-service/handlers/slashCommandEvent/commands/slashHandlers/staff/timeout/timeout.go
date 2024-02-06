package timeoutHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeout(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()
	reason := i.ApplicationCommandData().Options[1].StringValue()
	sTimeLengthString := i.ApplicationCommandData().Options[2].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeout` command..."),
		},
	})

	sTimeLength, convErr := utils.StringToFloat64(sTimeLengthString)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `timeLength` command argument is invalid. (term: %s)", sTimeLengthString)
		utils.SendErrorReportEmbed(s, i.Interaction, errMsg)
		return
	}

	timestamp := time.Now().Unix()

	var err error
	err = member.GiveTimeoutToMemberWithId(s, i, targetUserId, reason, timestamp, *sTimeLength)
	if err != nil {
		errMsg := fmt.Sprintf("Error ocurred giving timeout to user with ID %s: %s", targetUserId, err)
		utils.SendErrorReportEmbed(s, i.Interaction, errMsg)
		return
	}

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while sending timeout embed response: %v", err)
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId)
		utils.SendErrorReportEmbed(s, i.Interaction, errMsg)
		return
	}
	// TODO send DM

	// Format timeout creation time
	var timeoutCreatedAt time.Time
	var timeoutCreatedAtString string
	timeoutCreatedAt = time.Unix(timestamp, 0).UTC()
	timeoutCreatedAtString = timeoutCreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")

	// Format timeout duration
	var dd, hr, mm, ss = utils.HumanReadableTimeLength(*sTimeLength)
	var timeoutLengthString string = fmt.Sprintf("%dd, %dh:%dm:%ds", dd, hr, mm, ss)

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Timeout given to `%s`", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Reason", reason, false).
		AddField("Duration", timeoutLengthString, false).
		AddField("Timestamp", timeoutCreatedAtString, false)

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
