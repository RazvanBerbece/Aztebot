package slashHandlers

import (
	"fmt"
	"log"
	"os"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMusicDisconnect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := disconnectFromMusicVoiceChannel(s, i)
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

func disconnectFromMusicVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {

	if !globals.AztemusicApp.IsJoined || globals.AztemusicApp.VoiceChannel == nil {
		return fmt.Sprintf("%s is not currently on any voice channel.", os.Getenv("APP_NAME")), nil
	}

	err := globals.AztemusicApp.VoiceChannel.Disconnect()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error occured while disconnecting %s from VC: %s", os.Getenv("APP_NAME"), err))
		return "", err
	}
	globals.AztemusicApp.VoiceChannel = nil
	globals.AztemusicApp.IsJoined = false
	logging.LogVCOperation(fmt.Sprintf("Disconnected from voice channel with ID %s", globals.AztemusicApp.VoiceChannel.ChannelID), "")

	return "I have successfully disconnected from the voice channel.", nil

}
