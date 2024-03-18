package startup

import (
	"fmt"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SendInformationEmbedsToTextChannels(s *discordgo.Session) {

	// For each available default message resource in local storage
	for id, details := range globalConfiguration.DefaultInformationChannels {
		hasMessage, err := utils.ChannelHasDefaultInformationMessage(s, id)
		if err != nil {
			fmt.Printf("Could not check for default message in channel %s (%s): %v", id, details, err)
			continue
		}
		if hasMessage {
			// Do not send this default message as it already exists
			continue
		} else {
			// Send associated default message to given text channel
			var embedText string
			var hasOwnEmbed bool
			var longEmbed *embed.Embed
			switch details {
			case "default":
				embedText = utils.GetTextFromFile("internal/handlers/remoteEvents/readyEvent/assets/defaultContent/default.txt")
			case "info-music":
				embedText = utils.GetTextFromFile("internal/handlers/remoteEvents/readyEvent/assets/defaultContent/music-info.txt")
			case "staff-rules":
				embedText = utils.GetTextFromFile("internal/handlers/remoteEvents/readyEvent/assets/defaultContent/staff-rules.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			case "server-rules":
				embedText = utils.GetTextFromFile("internal/handlers/remoteEvents/readyEvent/assets/defaultContent/server-rules.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			case "legends":
				embedText = utils.GetTextFromFile("internal/handlers/remoteEvents/readyEvent/assets/defaultContent/legends.txt")
				hasOwnEmbed = true
				longEmbed = utils.GetLongEmbedFromStaticData(embedText)
			}

			var messageEmbedToPost *discordgo.MessageEmbed
			if !hasOwnEmbed {
				messageEmbedToPost = embed.NewEmbed().
					SetTitle("🤖  Information Message").
					SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
					SetColor(000000).
					AddField("", embedText, false).
					MessageEmbed
			} else {
				messageEmbedToPost = longEmbed.MessageEmbed
			}

			_, err := s.ChannelMessageSendEmbed(id, messageEmbedToPost)
			if err != nil {
				log.Fatalf("An error occured while sending a default message (%s): %v", details, err)
				return
			}
		}
	}

}
