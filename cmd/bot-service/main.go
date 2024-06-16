package main

import (
	discordBotBaseModule "github.com/RazvanBerbece/Aztebot/internal/bot-service/base"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers"
)

func main() {

	// Configure the bot base with the key, handlers and intents
	bot := discordBotBaseModule.DiscordBotBase{}
	bot.ConfigureBase("aztebot")
	bot.AddHandlers(handlers.GetAztebotHandlersAsList())

	// Connect to the Discord servers
	bot.Connect()

	// Close connection
	bot.CloseConnection()

	// Cleanup used resources
	bot.Cleanup()

}
