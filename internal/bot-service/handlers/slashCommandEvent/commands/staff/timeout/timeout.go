package timeoutSlashHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/notifications"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeout(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	reason := i.ApplicationCommandData().Options[1].StringValue()
	sTimeLengthString := i.ApplicationCommandData().Options[2].StringValue()

	commandOwnerUserId := i.Member.User.ID

	// Input validation
	if !utils.IsValidDiscordUserId(targetUserId) {
		errMsg := fmt.Sprintf("The provided `user-id` command argument is invalid. (term: `%s`)", targetUserId)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if !utils.IsValidReasonMessage(reason) {
		errMsg := fmt.Sprintf("The provided `reason` command argument is invalid. (term: `%s`)", reason)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeout` command..."),
		},
	})

	user, err := s.User(targetUserId)
	if err != nil {
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
	}

	sTimeLength, convErr := utils.StringToFloat64(sTimeLengthString)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `duration` command argument is invalid. (term: `%s`)", sTimeLengthString)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	// Validate the timeout duration input to be in the allowed array of values
	allowedTimeoutExpirations := []float64{300, 600, 1800, 3600, 86400, 259200, 604800}
	if !utils.Float64InSlice(*sTimeLength, allowedTimeoutExpirations) {
		errMsg := fmt.Sprintf("The provided `duration` command argument is not an allowed value. (term `%s` not in { 300, 600, 1800, 3600, 86400, 259200, 604800 })", sTimeLengthString)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	timestamp := time.Now().Unix()

	err = member.GiveTimeoutToMemberWithId(s, globals.DiscordMainGuildId, targetUserId, reason, timestamp, *sTimeLength)
	if err != nil {
		errMsg := fmt.Sprintf("Error ocurred giving timeout to user with UID `%s`: `%s`", targetUserId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	// Format timeout creation time
	var timeoutCreatedAt time.Time
	var timeoutCreatedAtString string
	timeoutCreatedAt = time.Unix(timestamp, 0).UTC()
	timeoutCreatedAtString = timeoutCreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")

	// Build a DM to send to the target user detailing the timeout
	timeoutDm := fmt.Sprintf("You received a timeout for reason: `%s`\nat `%s`", reason, timeoutCreatedAtString)
	err = member.SendDirectMessageToMember(s, targetUserId, timeoutDm)
	if err != nil {
		fmt.Printf("An error ocurred while sending timeout embed response: %v", err)
		errMsg := fmt.Sprintf("An error ocurred while sending the timeout DM to the target user %s", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
	}

	// Format timeout duration
	var dd, hr, mm, ss = utils.HumanReadableTimeLength(*sTimeLength)
	var timeoutLengthString string = fmt.Sprintf("%dd, %dh:%dm:%ds", dd, hr, mm, ss)

	// Send notification to target channel to announce the timeout
	if channel, channelExists := globals.NotificationChannels["notif-timeout"]; channelExists {
		go sendTimeoutNotification(s, channel.ChannelId, targetUserId, reason, timeoutCreatedAtString, timeoutLengthString, commandOwnerUserId)
	}

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

func sendTimeoutNotification(s *discordgo.Session, channelId string, targetUserId string, reason string, timestamp string, duration string, commandOwnerUserId string) {

	// Get command owner discord name
	cmdOwner, err := s.User(commandOwnerUserId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving command owner with ID: %v", err)
	}

	fields := []discordgo.MessageEmbedField{
		{
			Name:   "By Staff Member",
			Value:  cmdOwner.Username,
			Inline: false,
		},
		{
			Name:   "Reason",
			Value:  reason,
			Inline: false,
		},
		{
			Name:   "Duration",
			Value:  duration,
			Inline: false,
		},
		{
			Name:   "Created At",
			Value:  timestamp,
			Inline: false,
		},
	}

	notificationTitle := fmt.Sprintf("`/timeout` given to User with UID `%s`", targetUserId)
	notifications.SendNotificationToTextChannel(s, channelId, notificationTitle, fields, true)

}
