package commands

import (
	slashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers"
	"github.com/bwmarrin/discordgo"
)

var AztebotSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Basic ping slash interaction for the AzteBot",
	},
	{
		Name:        "my_roles",
		Description: "Get a list of your assigned roles",
	},
	{
		Name:        "me",
		Description: "Get a summary of your profile details which are linked to the OTA guild",
	},
	{
		Name:        "sync",
		Description: "Syncs the user profile data (roles, etc.) with the OTA servers. This is a restricted command",
	},
	{
		Name:        "help",
		Description: "Get a help guide for the available AzteBot slash commands",
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":     slashHandlers.HandleSlashPingAztebot,
	"my_roles": slashHandlers.HandleSlashMyRoles,
	"me":       slashHandlers.HandleSlashMe,
	"help":     slashHandlers.HandleSlashAztebotHelp,
	"sync":     slashHandlers.HandleSlashSync,
}
