package slashCommandEvent

import (
	"log"

	"github.com/LxrdVixxeN/Aztebot/internal/bot-service/globals"
	commands "github.com/LxrdVixxeN/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/bwmarrin/discordgo"
)

func RegisterSlashCommands(s *discordgo.Session) error {

	globals.RegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.SlashCommands))
	for index, cmd := range commands.SlashCommands {
		_, err := s.ApplicationCommandCreate(globals.AppId, globals.DiscordGuildId, cmd)
		if err != nil {
			return err
		}
		globals.RegisteredCommands[index] = cmd
	}

	// Add slash command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handlerFunc, ok := commands.SlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})

	return nil
}

func CleanupSlashCommands(s *discordgo.Session) {
	for _, cmd := range globals.RegisteredCommands {
		err := s.ApplicationCommandDelete(globals.AppId, globals.DiscordGuildId, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete %s slash command: %v", cmd.Name, err)
		}
	}
}
