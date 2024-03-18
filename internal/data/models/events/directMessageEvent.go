package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

type DirectMessageEvent struct {
	UserId          string
	TargetChannelId string
	Type            string
	Title           *string
	Fields          []discordgo.MessageEmbedField
	Embed           *embed.Embed
	UseThumbnail    *bool
}
