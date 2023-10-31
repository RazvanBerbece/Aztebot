package commands

import (
	commandHandlers "github.com/LxrdVixxeN/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers"
	"github.com/bwmarrin/discordgo"
)

var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping_aztebot",
		Description: "Basic ping slash interaction for the Aztebot",
	},
	{
		Name:        "my_roles",
		Description: "Get a list of your assigned roles",
	},
}

var SlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping_aztebot": commandHandlers.HandleSlashPingAztebot,
	"my_roles":     commandHandlers.HandleSlashMyRoles,
}
