package globals

import (
	"os"

	"github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/player"
	appState "github.com/RazvanBerbece/Aztebot/pkg/shared/app"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// =============== ENVIRONMENT VARIABLES ===============
var _ = godotenv.Load(".env")

var Environment = os.Getenv("ENVIRONMENT") // staging / production

var DiscordAztebotToken = os.Getenv("DISCORD_AZTEBOT_TOKEN")
var DiscordAztebotAppId = os.Getenv("DISCORD_AZTEBOT_APP_ID")

var DiscordGuildId = os.Getenv("DISCORD_GUILD_ID")
var DiscordRadioChannelId = os.Getenv("DISCORD_RADIO_CHANNEL_ID")

var MySqlRootConnectionString = os.Getenv("DB_ROOT_CONNSTRING") // in MySQL format (i.e. `root_username:root_password@<unix/tcp>(host:port)/db_name`)

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============
var AztemusicApp *appState.AztemusicApp
var Player *player.Player

// =============== RUNTIME VARIABLES (SLASH COMMANDS) ===============
var AztebotRegisteredCommands []*discordgo.ApplicationCommand
var AztemusicRegisteredCommands []*discordgo.ApplicationCommand
