package handlers

import (
	messageEventHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/messageEvent"
	readyEventHandlers "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/readyEvent"
)

// Handler functions for the AzteBot.
func GetAztebotHandlersAsList() []interface{} {
	return []interface{}{
		// <---- On Ready ---->
		readyEventHandlers.Ready,
		// <---- On Message Created ---->
		messageEventHandlers.Ping, messageEventHandlers.SimpleMsgReply,
		// <---- On Reaction Added ---->
		// <---- On Reaction Removed ---->
		// <---- On New Join ---->
	}
}

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
