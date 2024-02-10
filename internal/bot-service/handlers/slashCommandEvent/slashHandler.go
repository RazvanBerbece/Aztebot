package slashCommandEvent

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	commands "github.com/RazvanBerbece/Aztebot/internal/bot-service/handlers/slashCommandEvent/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterSlashHandler(s *discordgo.Session) {

	// This handler runs on EVERY slash command registered with the main bot application
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		appData := i.ApplicationCommandData()

		// If a higher-up restricted command is being executed
		// and allowed roles are configured, only allow a user with one of these roles to execute an app command
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

		// If a staff command
		// and staff roles are configured, only allow a user with one of these roles to execute an app command
		if utils.StringInSlice(appData.Name, globals.StaffCommands) && len(globals.StaffRoles) > 0 {
			if i.Type == discordgo.InteractionApplicationCommand {
				// Check if the user has the allowed role
				hasAllowedRole := false
				for _, role := range i.Member.Roles {
					roleObj, err := s.State.Role(i.GuildID, role)
					if err != nil {
						log.Println("Error getting role:", err)
						return
					}
					if utils.StringInSlice(roleObj.Name, globals.StaffRoles) {
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

		err := globalsRepo.UserStatsRepository.IncrementSlashCommandsUsedForUser(i.Member.User.ID)
		if err != nil {
			fmt.Printf("Error ocurred while incrementing slash commands for user %s: %v", i.Member.User.ID, err)
		}

		err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(i.Member.User.ID)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", i.Member.User.ID, err)
		}
		err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(i.Member.User.ID, time.Now().Unix())
		if err != nil {
			fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", i.Member.User.ID, err)
		}

		if handlerFunc, ok := commands.AztebotSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})
}
