package coinSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashAddCoins(s *discordgo.Session, i *discordgo.InteractionCreate) {

	commandOwnerUserId := i.Member.User.ID

	// This is a staff command however restrict specifically anyone lower than "Developer" from using it
	ownerStaffRole, err := member.GetMemberStaffRole(commandOwnerUserId, globalConfiguration.StaffRoles)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't check for command owner permissions: %v", err)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if !utils.StringInSlice(ownerStaffRole.DisplayName, []string{"Developer", "Dominus", "Arhitect"}) {
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, "Command owner doesn't have the right permissions to use this command.")
		return
	}

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	coinsStr := i.ApplicationCommandData().Options[1].StringValue()

	user, err := globalRepositories.UsersRepository.GetUser(targetUserId)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
		return
	}

	fCoins, convErr := utils.StringToFloat64(coinsStr)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `coins` command argument is invalid. (term: `%s`)", coinsStr)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	low := 0.0
	high := 50000.0
	if *fCoins <= low || *fCoins > high {
		errMsg := fmt.Sprintf("The provided `coins` command argument is invalid. (term: `%s` outside bounds)", coinsStr)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	globalMessaging.CoinAwardsChannel <- events.CoinAwardEvent{
		GuildId:  i.GuildID,
		UserId:   targetUserId,
		Funds:    *fCoins,
		Activity: "MANUAL-AWARD",
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetColor(000000).
		SetAuthor("AzteBot Coin Manager", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		AddField(fmt.Sprintf("Awarded `🪙 %.2f` AzteCoins to `%s`.", *fCoins, user.DiscordTag), "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}
