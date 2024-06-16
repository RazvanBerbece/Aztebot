package timeoutHandlers

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashTimeout(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()
	reason := i.ApplicationCommandData().Options[1].StringValue()
	hTimeLengthString := i.ApplicationCommandData().Options[2].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Processing `/timeout` command..."),
		},
	})

	hTimeLength, convErr := utils.StringToInt64(hTimeLengthString)
	if convErr != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The provided `timeLength` command argument is invalid. (term: %s)", hTimeLengthString),
			},
		})
		return
	}

	timestamp := time.Now().Unix()

	var err error
	go func() {
		err := member.GiveTimeoutToMemberWithId(s, i, targetUserId, reason, timestamp, *hTimeLength)
		if err != nil {
			fmt.Printf("An error ocurred while giving timeout to user: %v\n", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("An error ocurred while giving timeout to user with ID %s.", targetUserId),
				},
			})
			return
		}
	}()

	user, err := s.User(targetUserId)
	if err != nil {
		fmt.Printf("An error ocurred while sending timeout embed response: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error ocurred while retrieving user with ID %s provided in the slash command.", targetUserId),
			},
		})
	}

	// Format timeout creation time
	var timeoutCreatedAt time.Time
	var timeoutCreatedAtString string
	timeoutCreatedAt = time.Unix(timestamp, 0).UTC()
	timeoutCreatedAtString = timeoutCreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")

	// Format timeout time length
	// var timeoutLength time.Time
	var timeoutLengthString string

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖⚠️   Timeout given to `%s`", user.Username)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000).
		AddField("Reason", reason, false).
		AddField("Duration", timeoutLengthString, false).
		AddField("Timestamp", timeoutCreatedAtString, false)

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}
