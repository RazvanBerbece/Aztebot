package globals

import (
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// =============== ENVIRONMENT VARIABLES ===============
var _ = godotenv.Load(".env")

var Environment = os.Getenv("ENVIRONMENT") // staging / production

var DiscordAztebotToken = os.Getenv("DISCORD_AZTEBOT_TOKEN")
var DiscordAztebotAppId = os.Getenv("DISCORD_AZTEBOT_APP_ID")

var DiscordMainGuildId = os.Getenv("DISCORD_MAIN_GUILD_ID")
var DiscordGuildIds = os.Getenv("DISCORD_GUILD_IDS")

var MySqlRootConnectionString = os.Getenv("DB_ROOT_CONNSTRING") // in MySQL format (i.e. `root_username:root_password@<unix/tcp>(host:port)/db_name`)

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============
var RestrictedCommands = strings.Split(os.Getenv("RESTRICTED_COMMANDS"), ",")
var AllowedRoles = strings.Split(os.Getenv("ALLOWED_ROLES"), ",")

var UserSyncInterval, UserSyncIntervalErr = strconv.Atoi(os.Getenv("USER_SYNC_INTERVAL")) // in seconds

// =============== RUNTIME VARIABLES (SLASH COMMANDS) ===============
var AztebotRegisteredCommands []*discordgo.ApplicationCommand
