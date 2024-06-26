package serverSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashArcadeLadder(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/arcade-ladder` command..."),
		},
	})

	entries, err := globalRepositories.ArcadeLadderRepository.GetArcadeLadder()
	if err != nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
		return
	}

	// Build response embed with a detailed arcade ladder view
	embedToSend := embed.NewEmbed().
		SetAuthor("AzteBot", "https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetTitle("👾🎮   The OTA Arcade Ladder").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST")

	if len(entries) == 0 {
		embedToSend.AddField("", "There are no entries in the arcade ladder at the moment.", false)
	} else {
		for idx, entry := range entries {
			user, err := globalRepositories.UsersRepository.GetUser(entry.UserId)
			if err != nil {
				utils.ErrorEmbedResponseEdit(s, i.Interaction, err.Error())
				return
			}
			embedToSend.AddField("", fmt.Sprintf("%d. `%s` - Won `%d` arcades", idx+1, user.DiscordTag, entry.Wins), false)
		}
	}

	paginationRow := embed.GetPaginationActionRowForEmbed(globalMessaging.PreviousPageOnEmbedEventId, globalMessaging.NextPageOnEmbedEventId)
	globalMessaging.ComplexResponsesChannel <- events.ComplexResponseEvent{
		Interaction:   i.Interaction,
		Embed:         embedToSend,
		PaginationRow: &paginationRow,
	}

}
