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
		Description: "Commands which instructs AzteRadio to join the designated Radio channel.",
	},
	{
		Name:        "azteradio_disconnect",
		Description: "Commands which instructs AzteRadio to disconnect from the designated Radio channel.",
	},
}

var AzteradioSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"azteradio_ping":       slashHandlers.HandleSlashPingAzteradio,
	"azteradio_join":       slashHandlers.HandleSlashRadioJoin,
	"azteradio_disconnect": slashHandlers.HandleSlashRadioDisconnect,
}
