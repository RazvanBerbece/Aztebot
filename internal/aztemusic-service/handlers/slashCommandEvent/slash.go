package slashCommandEvent

import (
	"log"
	"os"

	commands "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/bwmarrin/discordgo"
)

func RegisterAzteradioSlashCommands(s *discordgo.Session) error {

	globals.AztemusicRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztemusicSlashCommands))
	for index, cmd := range commands.AztemusicSlashCommands {
		_, err := s.ApplicationCommandCreate(os.Getenv("CLIENT_ID"), globals.DiscordGuildId, cmd)
		if err != nil {
			return err
		}
		globals.AztemusicRegisteredCommands[index] = cmd
	}

	// Add slash command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handlerFunc, ok := commands.AztemusicSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})

	return nil
}

func CleanupAzteradioSlashCommands(s *discordgo.Session) {
	for _, cmd := range globals.AztemusicRegisteredCommands {
		err := s.ApplicationCommandDelete(os.Getenv("CLIENT_ID"), globals.DiscordGuildId, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete %s slash command: %v", cmd.Name, err)
		}
	}
}
