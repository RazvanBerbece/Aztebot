package commands

import (
	gamesHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/games"
	profileHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/profile"
	serverHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/server"
	timeoutHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/staff/timeout"
	warningHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/staff/warning"
	utilHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/slashHandlers/utils"
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
		Name:        "warn_remove_oldest",
		Description: "Removes a user's oldest warning",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user the warning was given to",
				Required:    true,
			},
		},
	},
	{
		Name:        "warns",
		Description: "View a list of a members's warnings",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user who was given the warnings",
				Required:    true,
			},
		},
	},
	{
		Name:        "timeout",
		Description: "Timeout a user's acitivity (block text and voice channels, but allow `/timeout-appeal`).",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user the timeout is given to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for which the timeout was given (max. 500 characters)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "duration",
				Description: "The duration of the timeout in seconds (300, 600, 1800, 3600, 86400, 259200)",
				Required:    true,
			},
		},
	},
	{
		Name:        "timeouts",
		Description: "See a user's active and archived timeouts.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user to see the associated timeouts for",
				Required:    true,
			},
		},
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":               utilHandlers.HandleSlashPingAztebot,
	"my_roles":           profileHandlers.HandleSlashMyRoles,
	"me":                 profileHandlers.HandleSlashMe,
	"help":               serverHandlers.HandleSlashAztebotHelp,
	"sync":               profileHandlers.HandleSlashSync,
	"top":                serverHandlers.HandleSlashTop,
	"dice":               gamesHandlers.HandleSlashDice,
	"warn":               warningHandlers.HandleSlashWarn,
	"warn_remove_oldest": warningHandlers.HandleSlashWarnRemoveOldest,
	"warns":              warningHandlers.HandleSlashWarns,
	"timeout":            timeoutHandlers.HandleSlashTimeout,
	"timeouts":           profileHandlers.HandleSlashTimeouts,
}
