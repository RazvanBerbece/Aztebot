package actionEvent

import (
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	actionEventConfessionApproval "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/actionEvents/confess"
	"github.com/bwmarrin/discordgo"
)

func HandleMessageComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	eventCustomId := i.MessageComponentData().CustomID

	// Future button event handlers should be added here by custom ID
	switch eventCustomId {
	case globalMessaging.ConfessionApprovalEventId:
		actionEventConfessionApproval.HandleApproveConfession(s, i)
	case globalMessaging.ConfessionDisprovalEventId:
		actionEventConfessionApproval.HandleDeclineConfession(s, i)
	}
}
