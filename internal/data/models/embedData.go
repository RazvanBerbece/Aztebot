package dataModels

import "github.com/bwmarrin/discordgo"

type EmbedData struct {
	CurrentPage int
	FieldData   *[]*discordgo.MessageEmbedField
}
