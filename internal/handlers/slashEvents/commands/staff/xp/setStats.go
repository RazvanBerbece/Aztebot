package xpRateSettingSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashSetStats(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	messagesSent := i.ApplicationCommandData().Options[1].StringValue()
	slashUsed := i.ApplicationCommandData().Options[2].StringValue()
	reactionsReceived := i.ApplicationCommandData().Options[3].StringValue()
	timeVc := i.ApplicationCommandData().Options[4].StringValue()
	timeMusic := i.ApplicationCommandData().Options[5].StringValue()

	user, err := globalRepositories.UsersRepository.GetUser(targetUserId)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
		return
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   Updated Stats For `%s`", user.DiscordTag)).
		SetColor(000000)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func sendXpRateChangeNotification(channelId string, activityName string, multiplierName string) {

	// Build global XP rate change embed
	embed := embed.
		NewEmbed().
		SetTitle("🤖📣	Global XP Rate Change Announcement").
		SetColor(000000).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg")

	if multiplierName == "Default OTA Value" {
		embed.AddField("", fmt.Sprintf("`%s` are now worth the default amount of experience points.", activityName), false)
	} else {
		embed.AddField("", fmt.Sprintf("`%s` are now worth `%s` as many experience points !", activityName, multiplierName), false)
	}

	embed.AtTagEveryone()

	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: channelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           embed,
	}

}

func getArgumentDisplayNames(activityInput string, multiplierInput string) (string, string) {

	var activityName string
	var multiplierName string

	switch activityInput {
	case "msg_send":
		activityName = "Message Sends"
	case "react_recv":
		activityName = "Reactions Received"
	case "slash_use":
		activityName = "Slash Commands Used"
	case "spent_vc":
		activityName = "Time Spent in Voice Channels"
	case "spent_music":
		activityName = "Time Spent Listening to Music"
	}

	switch multiplierInput {
	case "def":
		multiplierName = "Default OTA Value"
	case "1.5":
		multiplierName = "1.5x"
	case "2.0":
		multiplierName = "2x"
	case "3.0":
		multiplierName = "3x"
	}

	return activityName, multiplierName
}
