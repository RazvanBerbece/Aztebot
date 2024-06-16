package slashHandlers

import (
	"fmt"

	streams "github.com/RazvanBerbece/Aztebot/internal/azteradio-service/stream"
	youtubeApi "github.com/RazvanBerbece/Aztebot/internal/azteradio-service/stream/platforms/apis/youtube"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashRadioPlay(s *discordgo.Session, i *discordgo.InteractionCreate) {

	platform := "youtube" // todo: make this dynamic based on command parameter ?
	err := playRadio(s, platform)
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

func playRadio(s *discordgo.Session, platform string) (err error) {

	// TODO: Clean this function up
	switch platform {
	case "youtube":
		youtubeApi := youtubeApi.YouTubeApi{
			ApiName: "youtube",
		}
		stream := streams.Stream{
			PlatformName: platform,
			PlatformApi:  youtubeApi,
		}
		stream.DownloadSongs()
		stream.Play(s)
		return nil
	}

	return nil
}
