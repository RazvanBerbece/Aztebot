package slashHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashRadioJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := joinRadioVoiceChannel(s, i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func joinRadioVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(globals.DiscordGuildId, globals.DiscordRadioChannelId, false, true)
	if err != nil {
		return "", err
	}
	globals.AzteradioApp.VoiceChannel = vc
	globals.AzteradioApp.IsJoined = true
	logging.LogVCOperation(fmt.Sprintf("Joined voice channel with ID %s", globals.DiscordRadioChannelId), "")

	// It's recommended to wait some time before playing anything on the channel
	time.Sleep(250 * time.Millisecond)

	return "I have successfully joined the voice channel !", nil

}
