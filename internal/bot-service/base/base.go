package base

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	aztemusicSlashCommands "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/handlers/slashCommandEvent"
	aztebotSlashCommands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
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

func (b *DiscordBotBase) ConfigureBaseWithToken(appName string, token string) {

	// Create session based on the required app
	b.appName = appName
	b.isConnected = false

	switch b.appName {
	case "aztebot":
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Fatal("Could not create an AzteBot session: ", err)
		}
		b.botSession = session
	case "aztemusic":
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Fatal("Could not create an AzteMusic session: ", err)
		}
		b.botSession = session
	case "azteradio":
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Fatal("Could not create an AzteRadio session: ", err)
		}
		b.botSession = session
	}

}

// Initiates the instance's botSession with a fully configured discordgo session (auth, handlers, intents).
func (b *DiscordBotBase) AddHandlers(handlers []interface{}) {

	// Register custom handlers as callbacks for various events
	for _, handler := range handlers {
		b.botSession.AddHandler(handler)
	}

	// Register intents to allow bot operations on the Discord server (read chats, write messages, react, DM, etc.)
	b.botSession.Identify.Intents = getBotIntents()

	// Register slash commands based on app type
	switch b.appName {
	case "aztebot":
		err := aztebotSlashCommands.RegisterAztebotSlashCommands(b.botSession)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteBot: ", err)
		}
	case "aztemusic":
		err := aztemusicSlashCommands.RegisterAzteradioSlashCommands(b.botSession)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteMusic bot: ", err)
		}
	case "azteradio":
		err := aztemusicSlashCommands.RegisterAzteradioSlashCommands(b.botSession)
		if err != nil {
			log.Fatal("Error registering slash commands for AzteRadio: ", err)
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
	switch b.appName {
	case "aztebot":
		aztebotSlashCommands.CleanupAztebotSlashCommands(b.botSession)
	case "aztemusic":
		aztemusicSlashCommands.CleanupAzteradioSlashCommands(b.botSession)
	case "azteradio":
		aztemusicSlashCommands.CleanupAzteradioSlashCommands(b.botSession)
	}
}

// Gets the available bot intents.
// TODO: Make these more granular depending on bot features
func getBotIntents() discordgo.Intent {
	intents := discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.PermissionManageMessages |
		discordgo.PermissionManageServer |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildVoiceStates
	return intents
}
