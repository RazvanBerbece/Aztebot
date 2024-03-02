package handlers

import (
	joinEvent "github.com/RazvanBerbece/Aztebot/internal/handlers/guildJoinEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/guildRemoveEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/memberUpdateEvent"
	messageDeleteEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/messageDeleteEvent"
	messageEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/messageEvent"
	reactionAddEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/reactionAddEvent"
	reactionRemoveEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/reactionRemoveEvent"
	readyEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/readyEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/voiceStateUpdateEvent"
)

// Handler functions for the AzteBot.
func GetAztebotHandlersAsList() []interface{} {
	return []interface{}{
		// <---- On Ready ---->
		readyEventHandlers.Ready,
		// <---- On Message Created ---->
		messageEventHandlers.Any, messageEventHandlers.Ping, messageEventHandlers.SimpleMsgReply,
		// <---- On Message Deleted ---->
		messageDeleteEventHandlers.MessageDelete,
		// <---- On Reaction Added ---->
		reactionAddEventHandlers.ReactionAdd,
		// <---- On Reaction Removed ---->
		reactionRemoveEventHandlers.ReactionRemove,
		// <---- On New Join ---->
		joinEvent.GuildJoin,
		// <---- On Member Leaving Guild ---->
		guildRemoveEvent.GuildRemove,
		// <---- On Voice State Update ---->
		voiceStateUpdateEvent.VoiceStateUpdate,
		// <---- On Member Update ---->
		memberUpdateEvent.MemberRoleUpdate,
	}
}
