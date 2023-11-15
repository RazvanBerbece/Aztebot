package commands

import (
	slashHandlers "github.com/RazvanBerbece/Aztebot/internal/azteradio-service/handlers/slashCommandEvent/commands/slashHandlers"
	"github.com/bwmarrin/discordgo"
)

var AzteradioSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "azteradio_ping",
		Description: "Basic ping slash interaction for the AzteRadio",
	},
	{
		Name:        "azteradio_join",
		Description: "Command which instructs AzteRadio to join the designated Radio channel.",
	},
	{
		Name:        "azteradio_disconnect",
		Description: "Command which instructs AzteRadio to disconnect from the designated Radio channel.",
	},
	{
		Name:        "azteradio_play",
		Description: "Command which instructs AzteRadio to start playing music to the designated Radio channel.",
	},
}

var AzteradioSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"azteradio_ping":       slashHandlers.HandleSlashRadioPing,
	"azteradio_join":       slashHandlers.HandleSlashRadioJoin,
	"azteradio_disconnect": slashHandlers.HandleSlashRadioDisconnect,
	"azteradio_play":       slashHandlers.HandleSlashRadioPlay,
}
