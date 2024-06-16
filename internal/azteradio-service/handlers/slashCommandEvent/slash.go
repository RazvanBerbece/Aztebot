package slashCommandEvent

import (
	"log"

	commands "github.com/RazvanBerbece/Aztebot/internal/azteradio-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/bwmarrin/discordgo"
)

func RegisterAzteradioSlashCommands(s *discordgo.Session) error {

	globals.AzteradioRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AzteradioSlashCommands))
	for index, cmd := range commands.AzteradioSlashCommands {
		_, err := s.ApplicationCommandCreate(globals.DiscordAzteradioAppId, globals.DiscordGuildId, cmd)
		if err != nil {
			return err
		}
		globals.AzteradioRegisteredCommands[index] = cmd
	}

	// Add slash command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handlerFunc, ok := commands.AzteradioSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})

	return nil
}

func CleanupAzteradioSlashCommands(s *discordgo.Session) {
	for _, cmd := range globals.AzteradioRegisteredCommands {
		err := s.ApplicationCommandDelete(globals.DiscordAzteradioAppId, globals.DiscordGuildId, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete %s slash command: %v", cmd.Name, err)
		}
	}
}
