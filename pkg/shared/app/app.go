package app

import (
	"github.com/bwmarrin/discordgo"
)

type AzteradioApp struct {
	AppName      string
	BaseApp      interface{}
	VoiceChannel *discordgo.VoiceConnection
	IsJoined     bool
}

type AztebotApp struct {
	AppName string
	BaseApp interface{}
}
