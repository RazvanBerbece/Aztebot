package globals

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// =============== ENVIRONMENT VARIABLES ===============
var _ = godotenv.Load(".env") // Load

var Environment = os.Getenv("ENVIRONMENT") // staging / production

var DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
var AppId = os.Getenv("APP_ID")
var DiscordGuildId = os.Getenv("DISCORD_GUILD_ID")

var MySqlRootConnectionString = os.Getenv("DB_ROOT_CONNSTRING") // in format `root_username:root_password@tcp(host:port)/db_name-env_name`
var MySqlUserConnectionString = os.Getenv("DB_USER_CONNSTRING")

// =============== RUNTIME VARIABLES ===============
var RegisteredCommands []*discordgo.ApplicationCommand
