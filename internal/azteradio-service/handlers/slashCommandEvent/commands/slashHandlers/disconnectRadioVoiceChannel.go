package slashHandlers

import (
	"fmt"
	"log"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashRadioDisconnect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := disconnectFromRadioVoiceChannel(s, i)
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

func disconnectFromRadioVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {

	if !globals.AzteradioApp.IsJoined || globals.AzteradioApp.VoiceChannel == nil {
		return "AzteRadio is not currently on any voice channel.", nil
	}

	err := globals.AzteradioApp.VoiceChannel.Disconnect()
	if err != nil {
		log.Fatalf("Error occured while disconnecting AzteRadio from VC: %s", err)
		return "", err
	}
	globals.AzteradioApp.VoiceChannel = nil
	globals.AzteradioApp.IsJoined = false
	logging.LogVCOperation(fmt.Sprintf("Disconnected from voice channel with ID %s", globals.DiscordRadioChannelId), "")

	return "I have successfully disconnected from the voice channel.", nil

}
