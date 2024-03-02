package gamesSlashHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashDice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: GetDiceRollEmbed(i.Interaction.Member.User.Username),
		},
	})
}

// Returns an embed which contains a random dice roll.
func GetDiceRollEmbed(userDisplayName string) []*discordgo.MessageEmbed {

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Dice Roll", userDisplayName)).
		SetColor(000000)

	diceRoll := utils.GetRandomIntegerInRange(1, 6)

	embed.AddField(fmt.Sprintf("You rolled a `%d` 🎲", diceRoll), "", false)

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
