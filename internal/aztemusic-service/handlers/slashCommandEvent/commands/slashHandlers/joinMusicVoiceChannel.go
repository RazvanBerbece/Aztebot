package slashHandlers

import (
	"fmt"
	"log"
	"time"

	streams "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/stream"
	youtubeApi "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/stream/platforms/apis/youtube"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMusicJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := joinMusicVoiceChannel(s, i)
	if err != nil {
		log.Fatal(err)
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

func joinMusicVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {

	guildId := i.GuildID
	userId := i.Member.User.ID

	// Find the channel that the message came from.
	c, err := s.State.Channel(i.ChannelID)
	if err != nil {
		// Could not find channel.
		return "", err
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return "", err
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == userId {

			// Join the voice channel that the interaction creator is on.
			vc, err := s.ChannelVoiceJoin(guildId, vs.ChannelID, false, true)
			if err != nil {
				return "", err
			}
			logging.LogVCOperation(fmt.Sprintf("Joined voice channel with ID %s", vs.ChannelID), "")
			globals.AztemusicApp.VoiceChannel = vc
			globals.AztemusicApp.IsJoined = true

			// Init streaming capabilities
			youtubeApi := youtubeApi.YouTubeApi{
				ApiName: "youtube",
			}
			globals.AztemusicApp.Stream = streams.NewStream(youtubeApi.ApiName, youtubeApi, globals.AztemusicApp.VoiceChannel)
			logging.LogVCOperation("Initialised stream for VC", "")

			// It's recommended to wait some time before playing anything on the channel
			time.Sleep(250 * time.Millisecond)

			return "I have successfully joined the voice channel !", nil
		}
	}

	return "", fmt.Errorf("could not join voice channel")

}
