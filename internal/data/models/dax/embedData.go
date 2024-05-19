package dataModels

import (
	"github.com/bwmarrin/discordgo"
)

type EmbedData struct {
	ChannelId   string
	CurrentPage int
	FieldData   *[]*discordgo.MessageEmbedField
	Timestamp   float64
}
