package handlers

import (
	messageEventHandlers "github.com/LxrdVixxeN/Aztebot/internal/bot-service/handlers/messageEvent"
)

func GetHandlersAsList() []interface{} {
	// Add new handler methods here
	return []interface{}{
		// <---- On Ready ---->
		// <---- On Message Created ---->
		messageEventHandlers.Ping, messageEventHandlers.SimpleMsgReply,
		// <---- On Reaction Added ---->
		// <---- On Reaction Removed ---->
		// <---- On New Join ---->
	}
}
