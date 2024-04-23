package actionEventEmbedPagination

import (
	"fmt"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	actionEventsUtils "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/actionEvents/utils"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandlePaginateNextOnEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalPaginatedEmbedId := i.Message.ID
	originalPaginatedEmbedChannelId := i.Message.ChannelID

	// Respond to the button press (on help embed interaction source)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Paginating...",
			Flags:   1 << 6, // ephemeral response
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v\n", err)
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
	}
	utils.DeleteInteractionResponse(s, i.Interaction, 0)

	// Get original interaction if it can be found in the in-memory map
	embedData, exists := globalState.EmbedsToPaginate[originalPaginatedEmbedId]
	if !exists {
		fmt.Println("Failed to paginate (Next)")
		return
	} else {
		go actionEventsUtils.UpdatePaginatedEmbedPage(s, &embedData, "NEXT", originalPaginatedEmbedChannelId, originalPaginatedEmbedId, globalConfiguration.EmbedPageSize)
	}

}

func HandlePaginatePreviousOnEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalPaginatedEmbedId := i.Message.ID
	originalPaginatedEmbedChannelId := i.Message.ChannelID

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Paginating...",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		log.Printf("Error sending ACK message: %v", err)
		return
	}
	utils.DeleteInteractionResponse(s, i.Interaction, 0)

	// Get original interaction if it can be found in the in-memory map
	embedData, exists := globalState.EmbedsToPaginate[originalPaginatedEmbedId]
	if !exists {
		fmt.Println("Failed to paginate (Previous)")
		return
	} else {
		go actionEventsUtils.UpdatePaginatedEmbedPage(s, &embedData, "PREV", originalPaginatedEmbedChannelId, originalPaginatedEmbedId, globalConfiguration.EmbedPageSize)
	}

}
