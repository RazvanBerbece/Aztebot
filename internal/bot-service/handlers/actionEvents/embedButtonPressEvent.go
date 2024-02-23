package actionEvent

import (
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	actionEventConfessionApproval "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/actionEvents/confess"
	"github.com/bwmarrin/discordgo"
)

func HandleMessageComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		handleEmbedButtonPressEventHandler(s, i)
	}
}

func handleEmbedButtonPressEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	eventCustomId := i.MessageComponentData().CustomID

	// Future button event handlers should be added here by custom ID
	switch eventCustomId {
	case globals.ConfessionApprovalEventId:
		actionEventConfessionApproval.HandleApproveConfession(s, i)
	case globals.ConfessionDisprovalEventId:
		actionEventConfessionApproval.HandleDeclineConfession(s, i)
	}
}
