package arcadeLadderSlashHandlers

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashArcadeWinner(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	arcadeName := i.ApplicationCommandData().Options[1].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/arcade-winner` command..."),
		},
	})

	if globalAnnouncementsChannel, channelExists := globalConfiguration.NotificationChannels["notif-global"]; channelExists {
		err := member.GiveArcadeWin(targetUserId, arcadeName, globalAnnouncementsChannel.ChannelId)
		if err != nil {
			utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
			return
		}
	} else {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "The `/arcade-winner` command cannot be used without a designated notification channel to send arcade ladder related alerts to.")
		return
	}

	user, err := globalRepositories.UsersRepository.GetUser(targetUserId)
	if err != nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "an error ocurred while retrieving the Discord tag for the given member. the Win was recorded though so ignore this error.")
		return
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("👾🎮  Arcade Won by `%s` !", user.DiscordTag)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Edit for the final command response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
