package globals

import (
	"os"

	"github.com/joho/godotenv"
)

// =============== ENVIRONMENT VARIABLES ===============
var _ = godotenv.Load(".env") // Load

var Environment = os.Getenv("ENVIRONMENT") // staging / production

var DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
var AppId = os.Getenv("APP_ID")

var MySqlConnectionString = os.Getenv("DB_CONNSTRING") // in format `root_username:root_password@tcp(host:port)/db_name-env_name`

var GuildId = os.Getenv("GUILD_ID") // the id of the server the bot is in

// =============== MIGRATION RELATED VARIABLES ===============
var DbDriver = "mysql"
