package handlers

import (
	joinEvent "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/guildJoinEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/guildRemoveEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/memberUpdateEvent"
	messageDeleteEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/messageDeleteEvent"
	messageEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/messageEvent"
	reactionAddEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/reactionAddEvent"
	reactionRemoveEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/reactionRemoveEvent"
	readyEventHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/readyEvent"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/voiceStateUpdateEvent"
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
