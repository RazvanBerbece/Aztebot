package gamesSlashHandlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

var SizzlingSoundFilepath = "internal/handlers/slashEvents/commands/games/assets/audio/sizzling.dca"
var SizzlingSoundDataBuffer = make([][]byte, 0)

var CurrentlyPlayingAudio = false

func HandleSlashNewSizzlingSlot(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Load audio assets if necessary
	if len(SizzlingSoundDataBuffer) == 0 {
		err := loadSound(SizzlingSoundFilepath)
		if err != nil {
			utils.SendErrorEmbedResponse(s, i.Interaction, err.Error())
			return
		}
	}

	// Initial state
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: SlotEmbed(i.Interaction.Member.User.Username),
		},
	})

	// Sizzling...
	animationCount := 9
	voiceChannelID, _ := member.GetUserVoiceChannel(s, i.GuildID, i.Member.User.ID)
	if voiceChannelID != "" {
		// play sizzling sound if author is in a voice channel
		go AnimateSlotEmbed(s, *i.Interaction, animationCount)
		if !CurrentlyPlayingAudio {
			playSound(s, i.GuildID, voiceChannelID)
		}
	} else {
		go AnimateSlotEmbed(s, *i.Interaction, animationCount)
	}

}

func AnimateSlotEmbed(s *discordgo.Session, i discordgo.Interaction, animationCount int) {

	var frame []*discordgo.MessageEmbed
	for range animationCount {
		frame = SlotEmbed(i.Member.User.Username)
		editContent := ""
		editWebhook := discordgo.WebhookEdit{
			Content: &editContent,
			Embeds:  &frame,
		}
		s.InteractionResponseEdit(&i, &editWebhook)
	}

	// process final results (wins, etc) in frame above
	// TODO

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

	// Start speaking.
	vc.Speaking(true)

	CurrentlyPlayingAudio = true

	// Send the buffer data.
	for _, buff := range SizzlingSoundDataBuffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	CurrentlyPlayingAudio = false

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
