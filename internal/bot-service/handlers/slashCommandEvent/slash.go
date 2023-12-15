package slashCommandEvent

import (
	"fmt"
	"log"
	"strings"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	commands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterAztebotSlashCommands(s *discordgo.Session) error {

	// For each guild ID, register the commands
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

	// Add slash command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		appData := i.ApplicationCommandData()

		// If allowed roles are configured, only allow a user with one of these roles to execute an app command
		// The app commands which require role permissions are defined here
		if utils.StringInSlice(appData.Name, globals.RestrictedCommands) && len(globals.AllowedRoles) > 0 {
			if i.Type == discordgo.InteractionApplicationCommand {
				// Check if the user has the allowed role
				hasAllowedRole := false
				for _, role := range i.Member.Roles {
					roleObj, err := s.State.Role(i.GuildID, role)
					if err != nil {
						log.Println("Error getting role:", err)
						return
					}
					if utils.StringInSlice(roleObj.Name, globals.AllowedRoles) {
						hasAllowedRole = true
					}
					if hasAllowedRole {
						break
					}
				}

				if !hasAllowedRole {
					// If the user doesn't have the allowed role, send a response
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You do not have the required role to use this command.",
						},
					})
					return
				}
			}
		}

		userStatsRepo := repositories.NewUsersStatsRepository()
		err := userStatsRepo.IncrementSlashCommandsUsedForUser(i.Member.User.ID)
		if err != nil {
			fmt.Printf("Error ocurred while incrementing slash commands for user %s: %v", i.Member.User.ID, err)
		}

		if handlerFunc, ok := commands.AztebotSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})

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
