package actionEventConfessApproval

import (
	"fmt"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	actionEventsUtils "github.com/RazvanBerbece/Aztebot/internal/handlers/actionEvents/utils"
	supportSlashHandlers "github.com/RazvanBerbece/Aztebot/internal/handlers/slashCommandEvent/commands/support"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleApproveConfession(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalApprovalMessageId := i.Message.ID
	originalApprovalMessageChannelId := i.Message.ChannelID

	// Get original interaction if it can be found in the in-memory map
	confessionMessage, exists := globalState.ConfessionsToApprove[originalApprovalMessageId]
	if !exists {
		utils.SendErrorEmbedResponse(s, i.Interaction, "This confession message could not be found in the internal records.")
		return
	} else {
		// Send confession to designated text channel
		if channel, channelExists := globalConfiguration.NotificationChannels["notif-confess"]; channelExists {
			go supportSlashHandlers.SendApprovedConfessionNotification(s, channel.ChannelId, confessionMessage)
		}
		delete(globalState.ConfessionsToApprove, originalApprovalMessageChannelId)
	}

	// Respond to the button press
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You approved confession with ID `%s`. The message will now be forwarded to the designated channel", originalApprovalMessageId),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v\n", err)
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
	}

	// Cleanup
	go utils.DeleteInteractionResponse(s, i.Interaction, 3000)
	go actionEventsUtils.DisableButtonsForApprovalActionRow(s, originalApprovalMessageChannelId, originalApprovalMessageId, globalMessaging.ConfessionApprovalEventId, globalMessaging.ConfessionDisprovalEventId)
	go actionEventsUtils.UpdateApprovedActionRowOriginalMessage(s, i.Member.User.Username, "APPROVED", originalApprovalMessageChannelId, originalApprovalMessageId, globalMessaging.ConfessionApprovalEventId, globalMessaging.ConfessionDisprovalEventId)
}

func HandleDeclineConfession(s *discordgo.Session, i *discordgo.InteractionCreate) {

	originalApprovalMessageId := i.Message.ID
	originalApprovalMessageChannelId := i.Message.ChannelID

	// Get original interaction if it can be found in the in-memory map
	_, exists := globalState.ConfessionsToApprove[originalApprovalMessageId]
	if !exists {
		utils.SendErrorEmbedResponse(s, i.Interaction, "This confession message could not be found in the internal records.")
		return
	} else {
		delete(globalState.ConfessionsToApprove, originalApprovalMessageId)
	}

	// Respond to the button press
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You declined confession with ID `%s`", originalApprovalMessageId),
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v\n", err)
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
	}

	// Cleanup
	go utils.DeleteInteractionResponse(s, i.Interaction, 3000)
	go actionEventsUtils.DisableButtonsForApprovalActionRow(s, originalApprovalMessageChannelId, originalApprovalMessageId, globalMessaging.ConfessionApprovalEventId, globalMessaging.ConfessionDisprovalEventId)
	go actionEventsUtils.UpdateApprovedActionRowOriginalMessage(s, i.Member.User.Username, "DECLINED", originalApprovalMessageChannelId, originalApprovalMessageId, globalMessaging.ConfessionApprovalEventId, globalMessaging.ConfessionDisprovalEventId)
}
