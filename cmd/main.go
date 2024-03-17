package main

import (
	"github.com/RazvanBerbece/Aztebot/internal/base"
	discordHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents"
)

func main() {

	// Configure the bot base with the key, handlers and intents
	bot := base.DiscordBotBase{}
	bot.ConfigureBase()
	bot.AddHandlers(discordHandlers.GetAztebotHandlersAsList())

	// Cleanup used resources when program stops executing
	defer bot.Cleanup()

	// Connect to the Discord servers
	bot.Connect()

	// Close connection
	bot.CloseConnection()

}
