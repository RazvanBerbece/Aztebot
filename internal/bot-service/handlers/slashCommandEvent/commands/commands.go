package commands

import (
	slashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers"
	"github.com/bwmarrin/discordgo"
)

var AztebotSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping_aztebot",
		Description: "Basic ping slash interaction for the AzteBot",
	},
	{
		Name:        "my_roles",
		Description: "Get a list of your assigned roles",
	},
}

var AzteradioSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping_azteradio",
		Description: "Basic ping slash interaction for the AzteRadio",
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping_aztebot": slashHandlers.HandleSlashPingAztebot,
	"my_roles":     slashHandlers.HandleSlashMyRoles,
}
