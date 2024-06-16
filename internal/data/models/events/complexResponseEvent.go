package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

type ComplexResponseEvent struct {
	Interaction   *discordgo.Interaction
	Embed         *embed.Embed
	Title         *string
	Text          *string
	PaginationRow *discordgo.ActionsRow
}
