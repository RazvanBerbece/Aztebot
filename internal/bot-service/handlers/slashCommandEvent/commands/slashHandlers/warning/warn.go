package warning

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashWarn(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()
	reason := i.ApplicationCommandData().Options[1].StringValue()

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
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("An error ocurred while giving warning to user with ID %s.", targetUserId),
				},
			})
			return
		}
	}()

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while sending warning embed response: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId),
			},
		})
	}

	// Format CreatedAt
	var warnCreatedAt time.Time
	var warnCreatedAtString string
	warnCreatedAt = time.Unix(timestamp, 0).UTC()
	warnCreatedAtString = warnCreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")

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

	result := globalsRepo.WarnsRepository.GetWarningsCountForUser(userId)
	if result < 0 {
		fmt.Println("ERROR occured while getting all warnings count for user")
		return fmt.Errorf("ERROR SendWarnDmToUser")
	}

	// After-effects of warns - demotions, kicks, etc.
	switch result {
	case 0:
		// Send rule guide to user and tell them to follow it
		staffRules := utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/staff-rules.txt")
		dmContent := fmt.Sprintf("⚠️ You received a warning with reason: `%s`. You have %d out of 4 warnings.\nKeep in mind that on receiving 4 warnings you will be kicked out of the OTA community.\n\nSee below the OTA Staff rulebook.\n%s", reason, result+1, staffRules)
		err := sendWarnDmToUser(s, i, userId, dmContent)
		if err != nil {
			fmt.Printf("An error ocurred while sending staff rules DM to user: %v\n", err)
		}
	case 1:
		// 1 downgrade for role
		errDemote := member.DemoteMember(s, globals.DiscordMainGuildId, userId)
		if errDemote != nil {
			fmt.Printf("An error ocurred while demoting user: %v\n", errDemote)
			return errDemote
		}
		// Send demotion message
		demotionMessageContent := "⚠️ This is a message to inform you that you have been demoted from your Circle role as you received your second warning."
		err := sendWarnDmToUser(s, i, userId, demotionMessageContent)
		if err != nil {
			fmt.Printf("An error ocurred while sending demotion message content 1 DM to user: %v\n", err)
		}
	case 2:
		// 1 downgrade for role
		errDemote := member.DemoteMember(s, globals.DiscordMainGuildId, userId)
		if errDemote != nil {
			fmt.Printf("An error ocurred while demoting user: %v\n", errDemote)
			return errDemote
		}
		// Send demotion message
		demotionMessageContent := "⚠️ This is a message to inform you that you have been demoted from your Circle role as you received your third warning."
		err := sendWarnDmToUser(s, i, userId, demotionMessageContent)
		if err != nil {
			fmt.Printf("An error ocurred while sending demotion message content 2 DM to user: %v\n", err)
		}
	case 3:
		// Send demotion message
		kickMessageContent := "⚠️ This is a message to inform you that you have been kicked from the OTA community as you received your fourth, and final warning."
		err := sendWarnDmToUser(s, i, userId, kickMessageContent)
		if err != nil {
			fmt.Printf("An error ocurred while sending kick message content DM to user: %v\n", err)
		}
		// kick from guild, timeout
		err = member.KickMember(s, globals.DiscordMainGuildId, userId)
		if err != nil {
			fmt.Println("Error kicking member for receiving 4th warning:", err)
			return err
		}
	}

	err := globalsRepo.WarnsRepository.SaveWarn(userId, reason, timestamp)
	if err != nil {
		fmt.Printf("ERROR GiveWarnToUserWithId: %v", err)
		return err
	}

	return nil

}

func sendWarnDmToUser(s *discordgo.Session, i *discordgo.InteractionCreate, userId string, reason string) error {

	errDm := dm.SendHelpDmToUser(s, i, userId, reason)
	if errDm != nil {
		fmt.Println("Error sending DM: ", errDm)
		return errDm
	}
	return nil

}
