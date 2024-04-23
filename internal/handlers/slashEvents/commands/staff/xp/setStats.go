package xpSystemSlashHandlers

import (
	"fmt"

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

	messagesSentInt, convErr := utils.StringToInt(messagesSent)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `messages-sent` command argument is invalid. (term: `%s`)", messagesSent)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	slashUsedInt, convErr := utils.StringToInt(messagesSent)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `slash-cmd-used` command argument is invalid. (term: `%s`)", slashUsed)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	reactReceivedInt, convErr := utils.StringToInt(reactionsReceived)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `reactions-received` command argument is invalid. (term: `%s`)", reactionsReceived)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	timeVcFloat, convErr := utils.StringToFloat64(timeVc)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `time-vc` command argument is invalid. (term: `%s`)", timeVc)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	timeMusicFloat, convErr := utils.StringToFloat64(timeMusic)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `time-music` command argument is invalid. (term: `%s`)", timeMusic)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	err = globalRepositories.UserStatsRepository.SetStats(targetUserId, *messagesSentInt, *slashUsedInt, *reactReceivedInt, *timeVcFloat, *timeMusicFloat)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, fmt.Sprintf("an error ocurred while setting stats for user: %v", err))
		return
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle(fmt.Sprintf("🤖   Updated Stats For `%s`", user.DiscordTag)).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		SetColor(000000)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}
