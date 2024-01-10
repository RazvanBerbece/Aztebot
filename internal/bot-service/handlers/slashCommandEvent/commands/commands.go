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
		Description: "Syncs the user profile data (roles, etc.) with the OTA servers",
	},
	{
		Name:        "help",
		Description: "Get a help guide for the available AzteBot slash commands",
	},
	{
		Name:        "top",
		Description: "See the OTA leaderboard tops by activity category",
	},
	{
		Name:        "dice",
		Description: "Roll a 6-sided dice and try your luck",
	},
	{
		Name:        "warn",
		Description: "Gives a warning (with a provided reason message) to the user with the given ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user the warning is to be given to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for which the warning was given (max. 500 characters)",
				Required:    true,
			},
		},
	},
	{
		Name:        "warn_remove",
		Description: "Removes a warning with the provided ID from the user with the given ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "warn-id",
				Description: "The warning ID of the warning to remove from the user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user the warning was given to",
				Required:    true,
			},
		},
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":        slashHandlers.HandleSlashPingAztebot,
	"my_roles":    slashHandlers.HandleSlashMyRoles,
	"me":          slashHandlers.HandleSlashMe,
	"help":        slashHandlers.HandleSlashAztebotHelp,
	"sync":        slashHandlers.HandleSlashSync,
	"top":         slashHandlers.HandleSlashTop,
	"dice":        slashHandlers.HandleSlashDice,
	"warn":        slashHandlers.HandleSlashWarn,
	"warn_remove": slashHandlers.HandleSlashWarnRemove,
}
