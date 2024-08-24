package commands

import (
	gamesSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/games"
	profileSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/profile"
	serverSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/server"
	arcadeLadderSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/arcadeLadder"
	coinSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/coin"
	jailSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/jail"
	massPingSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/massPing"
	gainRatesSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/rates"
	repSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/rep"
	timeoutSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/timeout"
	warningSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/warning"
	xpSystemSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/staff/xp"
	supportSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/support"
	slashUtils "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands/utils"
	"github.com/bwmarrin/discordgo"
)

var AztebotSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Basic ping slash interaction for the AzteBot.",
	},
	{
		Name:        "my-roles",
		Description: "Get a list of your assigned roles.",
	},
	{
		Name:        "roles",
		Description: "See a user's role card.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User to see the role card for",
				Required:    true,
			},
		},
	},
	{
		Name:        "me",
		Description: "Get a summary of your profile details which are linked to the OTA guild.",
	},
	{
		Name:        "you",
		Description: "See a user's profile card",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User to see the profile card for",
				Required:    true,
			},
		},
	},
	{
		Name:        "help",
		Description: "Get a help guide for the available AzteBot slash commands.",
	},
	{
		Name:        "top5user",
		Description: "See the OTA leaderboard top 5s by activity category.",
	},
	{
		Name:        "dice",
		Description: "Roll a 6-sided dice and try your luck.",
	},
	{
		Name:        "sizzling",
		Description: "Generates a slot machine to use on a text channel.",
	},
	{
		Name:        "warn",
		Description: "Gives a warning (with a provided reason message) to the user with the given ID.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User the warning is given to",
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
		Name:        "warn-remove-oldest",
		Description: "Removes a user's oldest warning.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User the warning was given to",
				Required:    true,
			},
		},
	},
	{
		Name:        "warns",
		Description: "View a a member's warnings.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User who was given the warnings",
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
				Name:        "user",
				Description: "The Discord User ID the timeout is given to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for which the timeout was given (max. 500 characters)",
				Required:    true,
				MaxLength:   500,
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
				Name:        "user",
				Description: "The Discord User to see the associated timeouts for",
				Required:    true,
			},
		},
	},
	{
		Name:        "timeout-remove-active",
		Description: "Removes a user's current active timeout (and skip archiving it).",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The Discord User to remove the active timeout from",
				Required:    true,
			},
		},
	},
	{
		Name:        "timeout-appeal",
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
	{
		Name:        "top",
		Description: "Displays the global OTA leaderboard",
	},
	{
		Name:        "set-global-xp-rate",
		Description: "Sets the global XP gain rate for a specific activity.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "activity",
				Description: "Select the activity to set the XP gain rate for",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Message Sends",
						Value: "msg_send",
					},
					{
						Name:  "Reactions Received",
						Value: "react_recv",
					},
					{
						Name:  "Slash Commands Used",
						Value: "slash_use",
					},
					{
						Name:  "Time Spent in Voice Channels",
						Value: "spent_vc",
					},
					{
						Name:  "Time Spent Listening to Music",
						Value: "spent_music",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "multiplier",
				Description: "Select the gain rate multiplier for the specified activity",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Default",
						Value: "def",
					},
					{
						Name:  "1.5x",
						Value: "1.5",
					},
					{
						Name:  "2x",
						Value: "2.0",
					},
					{
						Name:  "3x",
						Value: "3.0",
					},
				},
			},
		},
	},
	{
		Name:        "set-gender",
		Description: "Sets your profile gender to the selected option.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "gender",
				Description: "Select the gender to set for your profile",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Male",
						Value: "male",
					},
					{
						Name:  "Female",
						Value: "female",
					},
					{
						Name:  "Nonbinary",
						Value: "nonbin",
					},
					{
						Name:  "Other",
						Value: "other",
					},
				},
			},
		},
	},
	{
		Name:        "jail",
		Description: "Jails the given user, giving them the designated 'jailed' role and removing their permissions.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to jail",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason to jail the user for",
				Required:    true,
				MaxLength:   500,
			},
		},
	},
	{
		Name:        "unjail",
		Description: "Unjails the given user, removing the designated 'jailed' role and returning their permissions.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to unjail",
				Required:    true,
			},
		},
	},
	{
		Name:        "jail-view",
		Description: "Displays a high level view of the OTA Jail.",
	},
	{
		Name:        "jailed-user",
		Description: "Retrieves a jailed user's OTA Jail record.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to see the record for",
				Required:    true,
			},
		},
	},
	{
		Name:        "monthly-leaderboard",
		Description: "Displays a high level view of the monthly activity leaderboard for the current month.",
	},
	{
		Name:        "daily-leaderboard",
		Description: "Displays a high level view of the daily activity leaderboard for the current day.",
	},
	{
		Name:        "arcade-ladder",
		Description: "Displays a high level view of the server's arcade ladder.",
	},
	{
		Name:        "arcade-winner",
		Description: "Assigns an arcade win to the given user. Also announces their win on the designated channel.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to assign an arcade win to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "arcade-name",
				Description: "The name of the arcade competition that the user won (e.g: Valorant)",
				Required:    true,
			},
		},
	},
	{
		Name:        "set-stats",
		Description: "Elevated privilege command to set a user's stats in the OTA records.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to set the stats for",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "messages-sent",
				Description: "How many messages sent to set for the user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "slash-cmd-used",
				Description: "How many slash commands used to set for the user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reactions-received",
				Description: "How many reactions received to set for the user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time-vc",
				Description: "How many *seconds* spent in total in voice channels to set for the user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time-music",
				Description: "How many *seconds* spent in total in music channels to set for the user",
				Required:    true,
			},
		},
	},
	{
		Name:        "add-coins",
		Description: "Adds AzteCoins points to a user's wallet.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to give XP to.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "coins",
				Description: "How much AzteCoin to award to a user.",
				Required:    true,
			},
		},
	},
	{
		Name:        "set-global-coin-rate",
		Description: "Sets the global coin gain rate for a specific activity.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "activity",
				Description: "Select the activity to set the coin gain rate for",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Message Sends",
						Value: "msg_send",
					},
					{
						Name:  "Reactions Received",
						Value: "react_recv",
					},
					{
						Name:  "Slash Commands Used",
						Value: "slash_use",
					},
					{
						Name:  "Time Spent in Voice Channels",
						Value: "spent_vc",
					},
					{
						Name:  "Time Spent Listening to Music",
						Value: "spent_music",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "multiplier",
				Description: "Select the gain rate multiplier for the specified activity",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Default",
						Value: "def",
					},
					{
						Name:  "1.5x",
						Value: "1.5",
					},
					{
						Name:  "2x",
						Value: "2.0",
					},
					{
						Name:  "3x",
						Value: "3.0",
					},
				},
			},
		},
	},
	{
		Name:        "reset-rep",
		Description: "Resets a user's rep score, setting it to 0.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to reset the stat for.",
				Required:    true,
			},
		},
	},
	{
		Name:        "mass-dm-announcement",
		Description: "Sends a mass DM to all the members in the server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "msg",
				Description: "The messages to send to all the members on the server.",
				Required:    true,
				MaxLength:   2048,
			},
		},
	},
	{
		Name:        "gain-rates",
		Description: "Displays the current values of the reward gain rates per activity on the server.",
	},
	{
		Name:        "reset-gain-rates",
		Description: "Resets all gain rates for this server.",
	},
}

var AztebotSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": slashUtils.HandleSlashPingAztebot,
	"help": serverSlashHandlers.HandleSlashAztebotHelp,

	"my-roles":   profileSlashHandlers.HandleSlashMyRoles,
	"roles":      profileSlashHandlers.HandleSlashYouRoles,
	"me":         profileSlashHandlers.HandleSlashMe,
	"you":        profileSlashHandlers.HandleSlashYou,
	"set-gender": profileSlashHandlers.HandleSlashSetGender,

	"dice":     gamesSlashHandlers.HandleSlashDice,
	"sizzling": gamesSlashHandlers.HandleSlashNewSizzlingSlot,
	"confess":  supportSlashHandlers.HandleSlashConfess,

	"top5user":      serverSlashHandlers.HandleSlashTop5Users,
	"top":           serverSlashHandlers.HandleSlashTop,
	"arcade-ladder": serverSlashHandlers.HandleSlashArcadeLadder,

	"daily-leaderboard":   serverSlashHandlers.HandleSlashDailyLeaderboard,
	"monthly-leaderboard": serverSlashHandlers.HandleSlashMonthlyLeaderboard,

	"add-coins":            coinSlashHandlers.HandleSlashAddCoins,
	"set-stats":            xpSystemSlashHandlers.HandleSlashSetStats,
	"reset-rep":            repSlashHandlers.HandleSlashResetRep,
	"arcade-winner":        arcadeLadderSlashHandlers.HandleSlashArcadeWinner,
	"mass-dm-announcement": massPingSlashHandlers.HandleSlashMassDm,

	"warn":               warningSlashHandlers.HandleSlashWarn,
	"warn-remove-oldest": warningSlashHandlers.HandleSlashWarnRemoveOldest,
	"warns":              warningSlashHandlers.HandleSlashWarns,

	"timeout":               timeoutSlashHandlers.HandleSlashTimeout,
	"timeouts":              timeoutSlashHandlers.HandleSlashTimeouts,
	"timeout-remove-active": timeoutSlashHandlers.HandleSlashTimeoutRemoveActive,
	"timeout-appeal":        timeoutSlashHandlers.HandleSlashTimeoutAppeal,

	"jail":        jailSlashHandlers.HandleSlashJail,
	"unjail":      jailSlashHandlers.HandleSlashUnjail,
	"jail-view":   jailSlashHandlers.HandleSlashJailView,
	"jailed-user": jailSlashHandlers.HandleSlashJailedUser,

	"gain-rates":           gainRatesSlashHandlers.HandleSlashServerGainRates,
	"set-global-coin-rate": gainRatesSlashHandlers.HandleSlashSetGlobalCoinRateForActivity,
	"set-global-xp-rate":   gainRatesSlashHandlers.HandleSlashSetGlobalXpRateForActivity,
	"reset-gain-rates":     gainRatesSlashHandlers.HandleSlashResetGainRates,
}
