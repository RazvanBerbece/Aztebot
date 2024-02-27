package globals

import (
	"os"
	"strconv"
	"strings"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
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

var NotificationChannelsPairs = strings.Split(os.Getenv("NOTIFICATION_CHANNELS"), ",")

var GlobalCommands = strings.Split(os.Getenv("GLOBAL_COMMANDS"), ",")

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============
var DefaultExperienceReward_MessageSent float64 = 0.5
var DefaultExperienceReward_SlashCommandUsed float64 = 0.45
var DefaultExperienceReward_ReactionReceived float64 = 0.33
var DefaultExperienceReward_InVc float64 = 0.005
var DefaultExperienceReward_InMusic float64 = 0.0025

var ExperienceReward_MessageSent float64 = 0.5
var ExperienceReward_SlashCommandUsed float64 = 0.45
var ExperienceReward_ReactionReceived float64 = 0.33
var ExperienceReward_InVc float64 = 0.005
var ExperienceReward_InMusic float64 = 0.0025

var NotificationChannels = make(map[string]dataModels.Channel)

var VoiceSessions = make(map[string]time.Time)
var StreamSessions = make(map[string]*time.Time)
var MusicSessions = make(map[string]map[string]*time.Time)
var DeafSessions = make(map[string]time.Time)

var LastUsedTopTimestamp = time.Now().Add(-60 * time.Minute)
var LastUsedTop5sTimestamp = time.Now().Add(-60 * time.Minute)

var ConfessionsToApprove = make(map[string]string)

// =============== RUNTIME GLOBAL CHANNELS ===============
var ExperienceGrantsChannel = make(chan dataModels.ExperienceGrant)

// =============== RUNTIME CUSTOM EVENT IDs (FOR BUTTON PRESS EVENT HANDLERS) ===============
var ConfessionApprovalEventId = "approve_confession"
var ConfessionDisprovalEventId = "decline_confession"

// =============== RUNTIME VARIABLES (SLASH COMMANDS) ===============
var AztebotRegisteredCommands []*discordgo.ApplicationCommand
