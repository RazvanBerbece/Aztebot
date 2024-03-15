package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/bwmarrin/discordgo"
)

type NotificationEvent struct {
	TargetChannelId string
	Type            string
	Title           *string
	Fields          []discordgo.MessageEmbedField
	Embed           *embed.Embed
	ActionRow       *discordgo.ActionsRow
	TextData        *string
	UseThumbnail    *bool
}
