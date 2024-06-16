package slashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMusicPlay(s *discordgo.Session, i *discordgo.InteractionCreate) {

	platform := "youtube" // todo: make this dynamic based on command parameter ?
	err := playEssentialsRadio(s, platform)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to play AzteRadio with platform <%s>", platform),
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Successfully playing AzteRadio with platform <%s>", platform),
		},
	})
}

func playEssentialsRadio(s *discordgo.Session, platform string) (err error) {

	// TODO: Clean this function up
	switch platform {
	case "youtube":
		globals.AztemusicApp.Stream.WithUrlsSourceFile().PlayFromLocalSource(s)
	}

	return nil
}
