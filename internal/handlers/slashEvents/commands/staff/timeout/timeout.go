package timeoutSlashHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
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
		errMsg := fmt.Sprintf("The provided `user` command argument is invalid. (term: `%s`)", targetUserId)
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
		return
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

	err = member.GiveTimeoutToMemberWithId(s, globalConfiguration.DiscordMainGuildId, targetUserId, reason, timestamp, *sTimeLength)
	if err != nil {
		errMsg := fmt.Sprintf("Error ocurred giving timeout to user with UID `%s`: `%s`", targetUserId, err)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	// Build a DM to send to the target user detailing the timeout
	creationTimestampString := utils.FormatUnixAsString(timestamp, "Mon, 02 Jan 2006 15:04:05 MST")
	timeoutDmTitle := ""
	timeoutDm := fmt.Sprintf("You received a timeout for reason: `%s`\nat `%s`", reason, creationTimestampString)

	// Publish DM event to announce timeout to target user
	globalMessaging.DirectMessagesChannel <- events.DirectMessageEvent{
		UserId: targetUserId,
		Title:  &timeoutDmTitle,
		Text:   &timeoutDm,
	}

	// Format timeout duration
	var dd, hr, mm, ss = utils.HumanReadableDuration(*sTimeLength)
	var timeoutLengthString string = fmt.Sprintf("%dd, %dh:%dm:%ds", dd, hr, mm, ss)

	// Send notification to target channel to announce the timeout
	if channel, channelExists := globalConfiguration.NotificationChannels["notif-timeout"]; channelExists {
		go sendTimeoutNotification(s, channel.ChannelId, targetUserId, reason, creationTimestampString, timeoutLengthString, commandOwnerUserId)
	}

	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖⚠️   Timeout given to `%s`", user.Username)).
		SetColor(000000).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		AddField("Reason", reason, false).
		AddField("Duration", timeoutLengthString, false).
		AddField("Timestamp", creationTimestampString, false)

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

	useThumbnail := true
	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Title:           &notificationTitle,
		Type:            "EMBED_WITH_TITLE_AND_FIELDS",
		Fields:          fields,
		UseThumbnail:    &useThumbnail,
	}

}
