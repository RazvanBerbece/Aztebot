package handlers

import (
	messageEventHandlers "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/handlers/messageEvent"
)

// Handler functions for the AzteRadio.
func GetAzteradioHandlersAsList() []interface{} {
	return []interface{}{
		// <---- On Ready ---->
		// <---- On Message Created ---->
		messageEventHandlers.Ping,
		// <---- On Reaction Added ---->
		// <---- On Reaction Removed ---->
	}
}
