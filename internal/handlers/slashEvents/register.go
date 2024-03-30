package slashCommandEvent

import (
	"fmt"
	"strings"
	"sync"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session, mainGuildOnly bool) error {

	fmt.Printf("[STARTUP] Registering %d slash commands...\n", len(commands.AztebotSlashCommands))

	// TODO: Optimise this one as it takes ~2 mins to finish executing and it seems to scale poorly with more slash commands
	err := RegisterGuildSlashCommands(s, globalConfiguration.DiscordAztebotAppId, mainGuildOnly, &globalConfiguration.DiscordMainGuildId)
	if err != nil {
		fmt.Printf("error in RegisterGuildSlashCommands: %v\n", err)
		return err
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

func RegisterGuildSlashCommands(s *discordgo.Session, appId string, mainGuildOnly bool, mainGuildId *string) error {

	var wg sync.WaitGroup
	var errGroup []string

	if mainGuildOnly {
		// Register commands only for the main guild
		// This is more performant when the bot is not supposed to be in more guilds
		globalState.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
		for index, cmd := range commands.AztebotSlashCommands {
			wg.Add(1)
			go func() {
				defer wg.Done()
				go SlashRegisterWorker(s, globalConfiguration.DiscordAztebotAppId, *mainGuildId, index, cmd, errGroup)
			}()
		}
	} else {
		// For each guild where the bot exists in, register the available commands
		guildIds := strings.Fields(globalConfiguration.DiscordGuildIds)
		for _, guildId := range guildIds {
			globalState.AztebotRegisteredCommands = make([]*discordgo.ApplicationCommand, len(commands.AztebotSlashCommands))
			for index, cmd := range commands.AztebotSlashCommands {
				wg.Add(1)
				go func() {
					defer wg.Done()
					go SlashRegisterWorker(s, globalConfiguration.DiscordAztebotAppId, guildId, index, cmd, errGroup)
				}()
			}
		}
	}

	wg.Wait()
	if len(errGroup) > 0 {
		// errors ocurred in registrations
		return fmt.Errorf("%d errors ocurred in RegisterGuildSlashCommands\nfor instance: %s", len(errGroup), errGroup[0])
	}

	return nil
}

func SlashRegisterWorker(s *discordgo.Session, appId string, guildId string, index int, cmd *discordgo.ApplicationCommand, errGroup []string) {
	_, err := s.ApplicationCommandCreate(appId, guildId, cmd)
	if err != nil {
		errString := fmt.Sprintf("an error ocurred while registering slash command %s: %v", cmd.Name, err)
		errGroup = append(errGroup, errString)
	}
	globalState.AztebotRegisteredCommands[index] = cmd
}
