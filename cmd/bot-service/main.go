package main

import (
	discordBotBaseModule "github.com/LxrdVixxeN/azteca-discord/internal/bot-service/base"
	"github.com/LxrdVixxeN/azteca-discord/internal/bot-service/handlers"
)

func main() {

	// Configure handler functions
	handlers := handlers.GetHandlersAsList()

	// Configure the bot base with the key, intents and handlers
	bot := discordBotBaseModule.DiscordBotBase{}
	bot.Configure(handlers)

	// Connect to the Discord servers
	bot.Connect()

	// Close connection
	bot.CloseConnection()

	// Cleanup used resources
	bot.Cleanup()

}
