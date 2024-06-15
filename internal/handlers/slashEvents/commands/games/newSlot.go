package gamesSlashHandlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

var SizzlingSoundFilepath = "internal/handlers/slashEvents/commands/games/assets/audio/sizzling.dca"
var SizzlingSoundDataBuffer = make([][]byte, 0)

func HandleSlashNewSizzlingSlot(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Load audio assets
	err := loadSound(SizzlingSoundFilepath)
	if err != nil {
		utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
		return
	}

	// Initial state
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: SlotEmbed(i.Interaction.Member.User.Username),
		},
	})

	// Sizzling...
	voiceChannelID, _ := member.GetUserVoiceChannel(s, i.GuildID, i.Member.User.ID)
	if voiceChannelID != "" {
		// PLAY SIZZLING SOUND
		go playSound(s, i.GuildID, voiceChannelID)
	}

	animationCount := 11
	for range animationCount {
		final := SlotEmbed(i.Interaction.Member.User.Username)
		editContent := ""
		editWebhook := discordgo.WebhookEdit{
			Content: &editContent,
			Embeds:  &final,
		}
		s.InteractionResponseEdit(i.Interaction, &editWebhook)
	}

	// Final state & processing
	final := SlotEmbed(i.Interaction.Member.User.Username)
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &final,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)

}

func SlotEmbed(userDisplayName string) []*discordgo.MessageEmbed {

	fruitEmojis := []string{"🍉", "🍎", "🍋", "🍇", "🍊", "🍍", "7️⃣"}

	defaultSlotStateString := &strings.Builder{}
	defaultSlotState := [][]string{
		{utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis)},
		{utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis)},
		{utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis)},
		{utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis)},
		{utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis), utils.GetRandomFromArray(fruitEmojis)},
	}

	table := tablewriter.NewWriter(defaultSlotStateString)
	for _, v := range defaultSlotState {
		table.Append(v)
	}
	table.SetCenterSeparator("-")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Render()

	embed := embed.NewEmbed().
		SetTitle("🎰	Sizzling Hot").
		SetColor(000000).
		SetDescription(defaultSlotStateString.String())

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(filepath string) error {

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		SizzlingSoundDataBuffer = append(SizzlingSoundDataBuffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range SizzlingSoundDataBuffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
