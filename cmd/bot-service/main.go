package main

import (
	discordBotBaseModule "github.com/LxrdVixxeN/Aztebot/internal/bot-service/base"
	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/handlers"
)

func main() {

	// Configure the bot base with the key, intents and handlers
	bot := discordBotBaseModule.DiscordBotBase{}
	bot.Configure(handlers.GetHandlersAsList())

	// Connect to the Discord servers
	bot.Connect()

	// Close connection
	bot.CloseConnection()

	// Cleanup used resources
	bot.Cleanup()

}
