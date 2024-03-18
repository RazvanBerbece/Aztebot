package warningSlashHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashWarn(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := utils.GetDiscordIdFromMentionFormat(i.ApplicationCommandData().Options[0].StringValue())
	reason := i.ApplicationCommandData().Options[1].StringValue()

	if !utils.IsValidDiscordUserId(targetUserId) {
		errMsg := fmt.Sprintf("The provided `user` command argument is invalid. (term: `%s`)", targetUserId)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}
	if !utils.IsValidReasonMessage(reason) {
		errMsg := fmt.Sprintf("The provided `reason` command argument is invalid. (term: `%s`)", reason)
		utils.SendCommandErrorEmbedResponse(s, i.Interaction, errMsg)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/warn` command..."),
		},
	})

	timestamp := time.Now().Unix()
	var err error
	go func() {
		err := GiveWarnToUserWithId(s, i, targetUserId, reason, timestamp)
		if err != nil {
			fmt.Printf("An error ocurred while giving warning to user: %v\n", err)
			errMsg := fmt.Sprintf("An error ocurred while giving warning to user with ID `%s`.", targetUserId)
			utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
			return
		}
	}()

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while sending warning embed response: %v", err)
		errMsg := fmt.Sprintf("An error ocurred while retrieving user with ID `%s` provided in the slash command.", targetUserId)
		utils.ErrorEmbedResponseEdit(s, i.Interaction, errMsg)
		return
	}

	// Format CreatedAt
	warnCreatedAtString := utils.FormatUnixAsString(timestamp, "Mon, 02 Jan 2006 15:04:05 MST")

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Warning given to `%s`", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Reason", reason, false).
		AddField("Timestamp", warnCreatedAtString, false)

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}

func GiveWarnToUserWithId(s *discordgo.Session, i *discordgo.InteractionCreate, userId string, reason string, timestamp int64) error {

	result := globalRepositories.WarnsRepository.GetWarningsCountForUser(userId)
	if result < 0 {
		fmt.Println("ERROR occured while getting all warnings count for user")
		return fmt.Errorf("ERROR publishWarnDmForMember")
	}

	// After-effects of warns - demotions, kicks, etc.
	switch result {
	case 0:
		// Send rule guide to user and tell them to follow it
		staffRules := utils.GetTextFromFile("internal/handlers/readyEvent/assets/defaultContent/staff-rules.txt")
		dmContent := fmt.Sprintf("⚠️ You received a warning with reason: `%s`. You have %d out of 4 warnings.\nKeep in mind that on receiving 4 warnings you will be kicked out of the OTA community.\n\nSee below the OTA Staff rulebook.\n%s", reason, result+1, staffRules)
		go publishWarnDmForMember(userId, dmContent)
	case 1:
		// 1 downgrade for staff role
		demoteType := "STAFF"
		errDemote := member.DemoteMember(s, globalConfiguration.DiscordMainGuildId, userId, demoteType)
		if errDemote != nil {
			fmt.Printf("An error ocurred while demoting user: %v\n", errDemote)
			return errDemote
		}
		// Send demotion message
		demotionMessageContent := fmt.Sprintf("⚠️ This is a message to inform you that you have been demoted from your %s role as you received your second warning.", demoteType)
		go publishWarnDmForMember(userId, demotionMessageContent)
	case 2:
		// 1 downgrade for role
		demoteType := "STAFF"
		errDemote := member.DemoteMember(s, globalConfiguration.DiscordMainGuildId, userId, "STAFF")
		if errDemote != nil {
			fmt.Printf("An error ocurred while demoting user: %v\n", errDemote)
			return errDemote
		}
		// Send demotion message
		demotionMessageContent := fmt.Sprintf("⚠️ This is a message to inform you that you have been demoted from your %s role as you received your third warning.", demoteType)
		go publishWarnDmForMember(userId, demotionMessageContent)
	case 3:
		// Send kick message
		kickMessageContent := "⚠️ This is a message to inform you that you have been kicked from the OTA community as you received your fourth, and final warning."
		go publishWarnDmForMember(userId, kickMessageContent)
		// kick from guild, timeout
		err := member.KickMember(s, globalConfiguration.DiscordMainGuildId, userId)
		if err != nil {
			fmt.Println("Error kicking member for receiving 4th warning:", err)
			return err
		}
	}

	err := globalRepositories.WarnsRepository.SaveWarn(userId, reason, timestamp)
	if err != nil {
		fmt.Printf("ERROR GiveWarnToUserWithId: %v", err)
		return err
	}

	return nil

}

func publishWarnDmForMember(userId string, reason string) {

	warnTitle := ""
	globalMessaging.DirectMessagesChannel <- events.DirectMessageEvent{
		UserId: userId,
		Title:  &warnTitle,
		Text:   &reason,
	}

}
