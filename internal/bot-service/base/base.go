package base

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"

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

	// Register intents to allow bot operations on the Discord server (read chats, write messages, react, DM, etc.)
	b.setBotIntents()

	// Register slash commands based on app type
	err := aztebotSlashCommands.RegisterAztebotSlashCommands(b.botSession)
	if err != nil {
		log.Fatal("Error registering slash commands for AzteBot: ", err)
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
	fmt.Println("Discord bot session is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

// Closes the existing persistent websocket connection to the Discord servers.
func (b *DiscordBotBase) CloseConnection() {
	b.botSession.Close()
	b.isConnected = false
}

// Cleans up any used resources by the bot service.
func (b *DiscordBotBase) Cleanup() {
	// Cleanup resources
	aztebotSlashCommands.CleanupAztebotSlashCommands(b.botSession)
}

// Sets the required bot intents.
// TODO: Make these more granular depending on bot features
func (b *DiscordBotBase) setBotIntents() {
	b.botSession.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.PermissionManageMessages |
		discordgo.PermissionManageServer |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildVoiceStates
}

func (b *DiscordBotBase) setBotStateTrackers() {
	b.botSession.State.TrackVoice = true
	b.botSession.State.MaxMessageCount = 250
}
