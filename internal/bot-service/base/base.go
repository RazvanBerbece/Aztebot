package botbase

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LxrdVixxeN/azteca-discord/internal/bot-service/globals"
	"github.com/bwmarrin/discordgo"
)

type DiscordBotBase struct {
	botSession  *discordgo.Session
	isConnected bool
}

// Initiates the instance's botSession with a fully configured discordgo session (auth, handlers, intents).
func (b *DiscordBotBase) Configure(handlers []interface{}) {

	// Create session
	session, err := discordgo.New("Bot " + globals.DiscordBotToken)
	if err != nil {
		log.Fatal("Could not create a Discord Bot session. Err: ", err)
	}

	// Register custom handlers as callbacks for various events
	for _, handler := range handlers {
		session.AddHandler(handler)
	}

	// Register intents
	// TODO: Make these more granular depending on bot features
	intents := discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.PermissionManageMessages |
		discordgo.PermissionManageServer |
		discordgo.IntentsDirectMessages
	session.Identify.Intents = intents

	// Register slash commands
	// TODO

	b.botSession = session

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
	// TODO
}
