package actionEventEmbedPagination

import (
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandlePaginateNextOnEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalPaginatedEmbedId := i.Message.ID
	originalPaginatedEmbedChannelId := i.Message.ChannelID

	// Get original interaction if it can be found in the in-memory map
	embedDataToPaginate, exists := globalState.EmbedsToPaginate[originalPaginatedEmbedId]
	if !exists {
		utils.SendErrorEmbedResponse(s, i.Interaction, "This paginated embed could not be found in the internal bot state.")
		return
	} else {
		// Next page on the embed
		// TODO
	}

}

func HandlePaginatePreviousOnEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {
}
