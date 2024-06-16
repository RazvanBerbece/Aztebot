package slashCommandEvent

import (
	"log"
	"strings"

	commands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session) error {

	// For each guild ID, register the commands
	guildIds := strings.Fields(globals.DiscordGuildId)
	for _, guildId := range guildIds {
		globals.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
		for index, cmd := range commands.AztebotSlashCommands {
			_, err := s.ApplicationCommandCreate(globals.DiscordAztebotAppId, guildId, cmd)
			if err != nil {
				return err
			}
			globals.AztebotRegisteredCommands[index] = cmd
		}
	}

	// Add slash command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handlerFunc, ok := commands.AztebotSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})

	return nil
}

func CleanupAztebotSlashCommands(s *discordgo.Session) {
	// For each guild ID, cleanup the commands
	guildIds := strings.Fields(globals.DiscordGuildId)
	for _, guildId := range guildIds {
		for _, cmd := range globals.AztebotRegisteredCommands {
			err := s.ApplicationCommandDelete(globals.DiscordAztebotAppId, guildId, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %s slash command: %v", cmd.Name, err)
			}
		}
	}
}
