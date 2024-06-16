package base

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"

	aztebotSlashCommands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent"

	"github.com/bwmarrin/discordgo"
)

// A base API which integrates a Discord bot session and various helper methods to setup any specific bot application.
type DiscordBotBase struct {
	botSession  *discordgo.Session
	isConnected bool
	appName     string
}

func (b *DiscordBotBase) ConfigureBase(appName string) {

	// Create session based on the required app
	b.appName = appName
	b.isConnected = false

	session, err := discordgo.New("Bot " + globals.DiscordAztebotToken)
	if err != nil {
		log.Fatal("Could not create an AzteBot session: ", err)
	}
	b.botSession = session

}

// Initiates the instance's botSession with a fully configured discordgo session (auth, handlers, intents).
func (b *DiscordBotBase) AddHandlers(handlers []interface{}) {

	// Register custom handlers as callbacks for various events
	for _, handler := range handlers {
		b.botSession.AddHandler(handler)
	}

	// Allow specific state trackers
	b.setBotStateTrackers()

	// Register intents and permissions
	// to allow bot operations on the Discord server (read chats, write messages, react, DM, etc.)
	b.setBotPermissions()
	b.setBotIntents()

	// Register slash commands based on app type
	if globals.DiscordMainGuildId != "" {
		// Register slash commands only for main guild
		err := aztebotSlashCommands.RegisterAztebotSlashCommands(b.botSession, true)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteBot: ", err)
		}
	} else {
		// Register slash commands for all guilds
		err := aztebotSlashCommands.RegisterAztebotSlashCommands(b.botSession, false)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteBot: ", err)
		}
	}

}

// Opens a persistent websocket connection to the Discord servers. Note that this method waits
// until CTRL-C or anther term signal is received.
func (b *DiscordBotBase) Connect() {

	err := b.botSession.Open()
	if err != nil {
		log.Fatal("Could not connect the bot to the Discord servers. Err: ", err)
	}

	b.isConnected = true

	// wait here until CTRL-C or anther term signal is received
	fmt.Println("Discord bot session is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

// Closes the existing persistent websocket connection to the Discord servers.
func (b *DiscordBotBase) CloseConnection() {
	b.botSession.Close()
	b.isConnected = false
}

// Cleans up any used resources by the bot service.
func (b *DiscordBotBase) Cleanup() {
	// Store in-process data (runtime maps, etc.)
	storeInProgressData()
	// Cleanup resources
	aztebotSlashCommands.CleanupAztebotSlashCommands(b.botSession)
}

// Sets the required bot intents and permissions.
// TODO: Make these more granular depending on bot features
func (b *DiscordBotBase) setBotIntents() {
	b.botSession.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuildMembers |
		discordgo.IntentGuildVoiceStates |
		discordgo.IntentsGuildBans |
		discordgo.IntentsAllWithoutPrivileged
}

func (b *DiscordBotBase) setBotPermissions() {
	b.botSession.Identify.Intents = discordgo.PermissionManageMessages |
		discordgo.PermissionManageServer |
		discordgo.PermissionManageRoles |
		discordgo.PermissionManageChannels |
		discordgo.PermissionModerateMembers |
		discordgo.PermissionAll
}

func (b *DiscordBotBase) setBotStateTrackers() {
	b.botSession.State.TrackVoice = true
	b.botSession.State.MaxMessageCount = 250
}

func storeInProgressData() {

	fmt.Println("Storing data at cleanup time")

	for uid, joinTime := range globals.VoiceSessions {
		duration := time.Since(joinTime)
		err := globalsRepo.UserStatsRepository.AddToTimeSpentInVoiceChannels(uid, int(duration.Seconds()))
		if err != nil {
			fmt.Printf("An error ocurred while adding time spent to voice channels for user with id %s: %v", uid, err)
		}
		delete(globals.VoiceSessions, uid)
	}

}
