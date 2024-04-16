package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

type DirectMessageEvent struct {
	UserId    string
	Embed     *embed.Embed
	Title     *string
	Text      *string
	ActionRow *discordgo.ActionsRow
}
