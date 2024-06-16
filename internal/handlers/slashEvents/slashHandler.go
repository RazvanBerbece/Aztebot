package slashCommandEvent

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	actionEvent "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/actionEvents"
	"github.com/RazvanBerbece/Aztebot/internal/handlers/slashEvents/commands"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func AddRegisteredSlashEventHandlers(s *discordgo.Session) {

	// This handler runs on EVERY slash command registered with the main bot application
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		if i.Type == discordgo.InteractionMessageComponent {
			// This configures button press event handlers for the bot
			// i.e pressing 'Accept' on a button on a generated embed and emitting the event
			actionEvent.HandleMessageComponentInteraction(s, i)
			return
		}

		appData := i.ApplicationCommandData()

		// If a higher-up restricted command is being executed
		// and allowed roles are configured, only allow a user with one of these roles to execute an app command
		if utils.StringInSlice(appData.Name, globalConfiguration.RestrictedCommands) && len(globalConfiguration.AllowedRoles) > 0 {
			if i.Type == discordgo.InteractionApplicationCommand {
				// Check if the user has the allowed role
				hasAllowedRole := false
				for _, role := range i.Member.Roles {
					roleObj, err := s.State.Role(i.GuildID, role)
					if err != nil {
						log.Println("Error getting role:", err)
						return
					}
					if utils.StringInSlice(roleObj.Name, globalConfiguration.AllowedRoles) {
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
		if utils.StringInSlice(appData.Name, globalConfiguration.StaffCommands) && len(globalConfiguration.StaffRoles) > 0 {
			if i.Type == discordgo.InteractionApplicationCommand {
				// Check if the user has the allowed role
				hasAllowedRole := false
				for _, role := range i.Member.Roles {
					roleObj, err := s.State.Role(i.GuildID, role)
					if err != nil {
						log.Println("Error getting role:", err)
						return
					}
					if utils.StringInSlice(roleObj.Name, globalConfiguration.StaffRoles) {
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

		ownerUserId := i.Member.User.ID

		err := globalRepositories.UserStatsRepository.IncrementSlashCommandsUsedForUser(ownerUserId)
		if err != nil {
			fmt.Printf("Error ocurred while incrementing slash commands for user %s: %v", ownerUserId, err)
		}

		err = globalRepositories.UserStatsRepository.IncrementActivitiesTodayForUser(ownerUserId)
		if err != nil {
			fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v", ownerUserId, err)
		}
		err = globalRepositories.UserStatsRepository.UpdateLastActiveTimestamp(ownerUserId, time.Now().Unix())
		if err != nil {
			fmt.Printf("An error ocurred while udpating user (%s) last timestamp: %v", ownerUserId, err)
		}

		// Publish experience grant message on the channel
		globalMessaging.ExperienceGrantsChannel <- events.ExperienceGrantEvent{
			UserId: ownerUserId,
			Points: globalConfiguration.ExperienceReward_SlashCommandUsed,
		}

		if handlerFunc, ok := commands.AztebotSlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
			handlerFunc(s, i)
		}
	})
}
