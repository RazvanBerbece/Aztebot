package app

import (
	streams "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/stream"
	"github.com/bwmarrin/discordgo"
)

type AztemusicApp struct {
	AppName      string
	BaseApp      interface{}
	VoiceChannel *discordgo.VoiceConnection
	IsJoined     bool
	Stream       streams.Stream
}

type AztebotApp struct {
	AppName string
	BaseApp interface{}
}
