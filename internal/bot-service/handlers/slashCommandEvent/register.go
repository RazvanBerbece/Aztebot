package slashCommandEvent

import (
	"fmt"
	"strings"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	commands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session, mainGuildOnly bool) error {

	fmt.Printf("[STARTUP] Registering %d slash commands...\n", len(commands.AztebotSlashCommands))

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

	// Register global commands (available in bot DMs as well)
	RegisterDmCommands(s, globals.GlobalCommands)

	// Register actual slash command handler
	go RegisterSlashHandler(s)

	return nil
}

// Given the comprehensive list of slash commands registered on the bot,
// and a list of command names which depict which commands can be used in the bot's DMs,
// register the DM commands as global commands, so users can leverage them in DMs.
func RegisterDmCommands(s *discordgo.Session, dmCommands []string) {
	for _, cmd := range commands.AztebotSlashCommands {
		if utils.StringInSlice(cmd.Name, dmCommands) {
			// If a command that can be used in DMs too
			_, err := s.ApplicationCommandCreate(globals.DiscordAztebotAppId, "", cmd)
			if err != nil {
				fmt.Println("An error ocurred while registering DM (global) commands")
			}
		}
	}
}
