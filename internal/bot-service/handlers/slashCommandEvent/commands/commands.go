package commands

import (
	gamesSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/games"
	profileSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/profile"
	serverSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/server"
	timeoutSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/staff/timeout"
	warningSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/staff/warning"
	supportSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/support"
	slashUtils "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands/utils"
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
		Name:        "roles",
		Description: "See a user's role card",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user to see the role card for",
				Required:    true,
			},
		},
	},
	{
		Name:        "me",
		Description: "Get a summary of your profile details which are linked to the OTA guild",
	},
	{
		Name:        "you",
		Description: "See a user's profile card",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user to see the profile card for",
				Required:    true,
			},
		},
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
		Description: "View a a member's warnings",
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
				Description: "Select the duration of the timeout",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "5 minutes",
						Value: "300",
					},
					{
						Name:  "10 minutes",
						Value: "600",
					},
					{
						Name:  "30 minutes",
						Value: "1800",
					},
					{
						Name:  "1 hour",
						Value: "3600",
					},
					{
						Name:  "1 day",
						Value: "86400",
					},
					{
						Name:  "3 days",
						Value: "259200",
					},
					{
						Name:  "1 week",
						Value: "604800",
					},
				},
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
	{
		Name:        "timeout_remove_active",
		Description: "Removes a user's current active timeout (and skip archiving it).",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-id",
				Description: "The Discord User ID of the user to remove the active timeout from",
				Required:    true,
			},
		},
	},
	{
		Name:        "timeout_appeal",
		Description: "Appeal your current active timeout (if you have one)",
	},
	{
		Name:        "confess",
		Description: "Sends an anonymised confessional message to the designated text channel.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "confession-message",
				Description: "The confession message to post",
				Required:    true,
			},
		},
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":                  slashUtils.HandleSlashPingAztebot,
	"my_roles":              profileSlashHandlers.HandleSlashMyRoles,
	"roles":                 profileSlashHandlers.HandleSlashYouRoles,
	"me":                    profileSlashHandlers.HandleSlashMe,
	"you":                   profileSlashHandlers.HandleSlashYou,
	"sync":                  profileSlashHandlers.HandleSlashSync,
	"dice":                  gamesSlashHandlers.HandleSlashDice,
	"help":                  serverSlashHandlers.HandleSlashAztebotHelp,
	"top":                   serverSlashHandlers.HandleSlashTop,
	"confess":               supportSlashHandlers.HandleSlashConfess,
	"warn":                  warningSlashHandlers.HandleSlashWarn,
	"warn_remove_oldest":    warningSlashHandlers.HandleSlashWarnRemoveOldest,
	"warns":                 warningSlashHandlers.HandleSlashWarns,
	"timeout":               timeoutSlashHandlers.HandleSlashTimeout,
	"timeouts":              timeoutSlashHandlers.HandleSlashTimeouts,
	"timeout_remove_active": timeoutSlashHandlers.HandleSlashTimeoutRemoveActive,
	"timeout_appeal":        timeoutSlashHandlers.HandleSlashTimeoutAppeal,
}
