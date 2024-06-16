package globals

import (
	"os"
	"strconv"
	"strings"
	"time"

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

var RestrictedCommands = strings.Split(os.Getenv("RESTRICTED_COMMANDS"), ",")
var AllowedRoles = strings.Split(os.Getenv("ALLOWED_ROLES"), ",")

var StaffCommands = strings.Split(os.Getenv("STAFF_COMMANDS"), ",")
var StaffRoles = strings.Split(os.Getenv("STAFF_ROLES"), ",")

var FavourableActivitiesThreshold, FavourableActivitiesThresholdErr = strconv.Atoi(os.Getenv("FAVOURABLE_ACTIVITIES_THRESHOLD"))

var UpdateVoiceStateFrequency, UpdateVoiceStateFrequencyErr = strconv.Atoi(os.Getenv("UPDATE_VOICE_STATE_DURATIONS_FREQUENCY"))

var TimeoutClearFrequency, TimeoutClearFrequencyErr = strconv.Atoi(os.Getenv("TIMEOUT_CLEAR_FREQUENCY"))

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============
var VoiceSessions = make(map[string]time.Time)
var StreamSessions = make(map[string]*time.Time)
var MusicSessions = make(map[string]map[string]*time.Time)
var DeafSessions = make(map[string]time.Time)
var LastUsedTopTimestamp = time.Now().Add(-60 * time.Minute)

// =============== RUNTIME VARIABLES (SLASH COMMANDS) ===============
var AztebotRegisteredCommands []*discordgo.ApplicationCommand
