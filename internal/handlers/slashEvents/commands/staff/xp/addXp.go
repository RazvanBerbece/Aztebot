package xpSystemSlashHandlers

import (
	"fmt"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashAddXp(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	xpStr := i.ApplicationCommandData().Options[1].StringValue()

	user, err := globalRepositories.UsersRepository.GetUser(targetUserId)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
		return
	}

	xpFloat, convErr := utils.StringToFloat64(xpStr)
	if convErr != nil {
		errMsg := fmt.Sprintf("The provided `xp` command argument is invalid. (term: `%s`)", xpStr)
		utils.SendErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	// Give XP to member
	_, err = member.AddExperienceToMember(targetUserId, *xpFloat)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
		return
	}

	// Send response embed
	embed := embed.NewEmbed().
		SetColor(000000).
		AddField(fmt.Sprintf("Added `%.2f` XP to `%s`.", *xpFloat, user.DiscordTag), "", false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}
