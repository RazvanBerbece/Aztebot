package readyEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	logging.LogHandlerCall("Ready", "")

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "Type /help")

	// Other setups

	// Cron func to sync users and their DB entity
	// TODO

}
