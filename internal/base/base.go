package base

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	slashCommandEvent "github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents"
	"github.com/bwmarrin/discordgo"
)

// A base API which integrates a Discord bot session and various helper methods to setup any specific bot application.
type DiscordBotBase struct {
	botSession  *discordgo.Session
	isConnected bool
}

func (b *DiscordBotBase) ConfigureBase() {

	// Create session based on the required app
	b.isConnected = false

	session, err := discordgo.New("Bot " + globalConfiguration.DiscordAztebotToken)
	if err != nil {
		log.Fatal("Could not create an AzteBot session: ", err)
	}
	b.botSession = session

}

// Initiates the instance's botSession with a fully configured discordgo session (auth, handlers, intents).
func (b *DiscordBotBase) AddHandlers(handlers []interface{}) {

	// Register custom handlers as callbacks for various events
	fmt.Printf("[STARTUP] Registering %d event handlers...\n", len(handlers))
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
	if globalConfiguration.DiscordMainGuildId != "" {
		// Register slash commands only for main guild
		err := slashCommandEvent.RegisterAztebotSlashCommands(b.botSession, true)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteBot: ", err)
		}
	} else {
		// Register slash commands for all guilds
		err := slashCommandEvent.RegisterAztebotSlashCommands(b.botSession, false)
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
	fmt.Println("[STARTUP] Discord bot session is now running.")
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
	// TODO: Cleanup resources ?
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
		discordgo.PermissionReadMessageHistory |
		discordgo.PermissionManageServer |
		discordgo.PermissionManageRoles |
		discordgo.PermissionManageChannels |
		discordgo.PermissionModerateMembers |
		discordgo.PermissionAll
}

func (b *DiscordBotBase) setBotStateTrackers() {
	b.botSession.StateEnabled = true
	b.botSession.State.TrackVoice = true
	b.botSession.State.TrackChannels = true
	b.botSession.State.MaxMessageCount = 500
}
