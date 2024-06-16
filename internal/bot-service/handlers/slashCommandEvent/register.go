package slashCommandEvent

import (
	"log"
	"strings"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	commands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session, mainGuildOnly bool) error {

	if mainGuildOnly {
		// Register commands only for the main guild
		// This is more performant when the bot is not supposed to be in more guilds
		globals.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
		for index, cmd := range commands.AztebotSlashCommands {
			_, err := s.ApplicationCommandCreate(globals.DiscordAztebotAppId, globals.DiscordMainGuildId, cmd)
			if err != nil {
				return err
			}
			globals.AztebotRegisteredCommands[index] = cmd
		}
	} else {
		// For each guild where the bot exists in, register the available commands
		guildIds := strings.Fields(globals.DiscordGuildIds)
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
	}

	// Register actual slash command handler
	go RegisterSlashHandler(s)

	return nil
}

func CleanupAztebotSlashCommands(s *discordgo.Session) {
	// For each guild ID, cleanup the commands
	guildIds := strings.Fields(globals.DiscordGuildIds)
	for _, guildId := range guildIds {
		for _, cmd := range globals.AztebotRegisteredCommands {
			err := s.ApplicationCommandDelete(globals.DiscordAztebotAppId, guildId, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %s slash command: %v", cmd.Name, err)
			}
		}
	}
}
