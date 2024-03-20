package member

import (
	"fmt"
	"log"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func IsBot(s *discordgo.Session, guildId string, userId string, debug bool) (*bool, error) {

	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		if debug {
			log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		}
		return nil, err
	}

	isBot := apiUser.Bot

	return &isBot, nil
}

// Removes all roles on the actual Discord member.
func RemoveAllDiscordUserRoles(s *discordgo.Session, guildId string, userId string) error {

	// Get the member's roles
	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return err
	}

	// Find all user's roles and delete them
	for _, roleID := range member.Roles {
		err = s.GuildMemberRoleRemove(guildId, userId, roleID)
		if err != nil {
			fmt.Printf("Error removing role with ID %s: %v\n", roleID, err)
			return err
		}
	}

	return nil

}

func RemoveDiscordRoleFromMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Remove the role by role ID from the Discord member
	err := s.GuildMemberRoleRemove(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Println("Error removing role from Discord member:\n", err)
		return err
	}

	return nil

}

func GiveDiscordRoleToMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Add the role by role ID to the Discord member
	err := s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Printf("Error giving role with name %s to Discord member: %v\n", roleName, err)
		return err
	}

	return nil

}

func AddRolesToDiscordUser(s *discordgo.Session, guildId string, userId string, roleIds []int) error {

	// For each role
	for _, roleId := range roleIds {
		role, err := globalRepositories.RolesRepository.GetRoleById(roleId)
		if err != nil {
			fmt.Printf("Error ocurred while adding DB roles to Discord member: %v\n", err)
			return err
		}
		// Get the role ID by display name from Discord
		discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, role.DisplayName)
		if discordRoleId != nil {
			// Add the role by role ID to the Discord member
			err = s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
			if err != nil {
				fmt.Printf("Error adding DB role with name %s to Discord member: %v\n", role.DisplayName, err)
				return err
			}

			// If a staff role, add a default 'STAFF' role as well
			if utils.StringInSlice(role.DisplayName, globalConfiguration.StaffRoles) {
				discordDefaultStaffRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, "STAFF")
				if discordDefaultStaffRoleId != nil {
					err = s.GuildMemberRoleAdd(guildId, userId, *discordDefaultStaffRoleId)
					if err != nil {
						fmt.Println("Error adding default STAFF role to Discord member:", err)
						return err
					}
				}
			}
		}
	}

	// Process ORDER role post-update (based on the current role state)
	_, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
	if currentOrder != nil {
		var discordOrderRoleIdToAdd *string
		if *currentOrder == 3 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
		} else if *currentOrder == 2 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
		} else if *currentOrder == 1 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
		}
		err := s.GuildMemberRoleAdd(guildId, userId, *discordOrderRoleIdToAdd)
		if err != nil {
			fmt.Println("Error adding order role to member:", err)
			return err
		}
	}

	return nil

}

func GetDiscordRoleIdForRoleWithName(s *discordgo.Session, guildId string, roleName string) *string {

	// Get the Guild
	guild, err := s.Guild(guildId)
	if err != nil {
		fmt.Println("Error retrieving guild:", err)
		return nil
	}

	// Find the Role ID based on the role's display name
	var roleID string = ""
	for _, role := range guild.Roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		fmt.Println("No role ID was found for role name", roleName)
		return nil
	}

	return &roleID
}
