package slashHandlers

import (
	"fmt"
	"time"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashWarn(s *discordgo.Session, i *discordgo.InteractionCreate) {

	targetUserId := i.ApplicationCommandData().Options[0].StringValue()
	reason := i.ApplicationCommandData().Options[1].StringValue()

	err := GiveWarnToUserWithId(targetUserId, reason, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while giving warning to user: %v\n", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error ocurred while giving warning to user with ID %s.", targetUserId),
			},
		})
	}

	// TODO: Make a nice embed to show the reason, timestamp and discord tag etc.
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Warning successfully given to user with ID %s.", targetUserId),
		},
	})

}

func GiveWarnToUserWithId(userId string, reason string, timestamp int64) error {

	err := globalsRepo.WarnsRepository.SaveWarn(userId, reason, timestamp)
	if err != nil {
		fmt.Printf("ERROR GiveWarnToUserWithId: %v", err)
		return err
	}

	return nil

}
