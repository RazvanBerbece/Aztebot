package slashCommandEvent

import (
	"fmt"
	"strings"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session, mainGuildOnly bool) error {

	fmt.Printf("[STARTUP] Registering %d slash commands...\n", len(commands.AztebotSlashCommands))

	// TODO: Optimise this as it take ~2 mins to finish executing and it seems to scale poorly with more slash commands

	if mainGuildOnly {
		// Register commands only for the main guild
		// This is more performant when the bot is not supposed to be in more guilds
		globalState.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
		for index, cmd := range commands.AztebotSlashCommands {
			_, err := s.ApplicationCommandCreate(globalConfiguration.DiscordAztebotAppId, globalConfiguration.DiscordMainGuildId, cmd)
			if err != nil {
				return err
			}
			globalState.AztebotRegisteredCommands[index] = cmd
		}
	} else {
		// For each guild where the bot exists in, register the available commands
		guildIds := strings.Fields(globalConfiguration.DiscordGuildIds)
		for _, guildId := range guildIds {
			globalState.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
			for index, cmd := range commands.AztebotSlashCommands {
				_, err := s.ApplicationCommandCreate(globalConfiguration.DiscordAztebotAppId, guildId, cmd)
				if err != nil {
					return err
				}
				globalState.AztebotRegisteredCommands[index] = cmd
			}
		}
	}

	// Register global commands (available in bot DMs as well)
	go RegisterDmCommands(s, globalConfiguration.GlobalCommands)

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
			_, err := s.ApplicationCommandCreate(globalConfiguration.DiscordAztebotAppId, "", cmd)
			if err != nil {
				fmt.Println("An error ocurred while registering DM (global) commands")
			}
		}
	}
}
