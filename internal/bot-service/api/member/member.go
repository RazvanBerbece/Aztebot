package member

import (
	"fmt"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func DemoteMember(s *discordgo.Session, guildId string, userId string) error {

	userToUpdate, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, err)
		return err
	}

	// DEMOTE STRATEGY (NOTH INNER CIRCLE ORDERS AND STAFF ROLE DEMOTIONS)
	userRoles, errUsrRole := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if errUsrRole != nil {
		fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, errUsrRole)
		return err
	}

	// DEMOTE IN THE DB
	var updatedCurrentRoleIds string = ""
	var roleIdsPriorDemote []int
	var roleIdsPostDemote []int
	var roleBeforeDemotion dataModels.Role
	for _, role := range userRoles {
		// If an Inner Circle role
		if role.Id > 7 && role.Id < 18 {
			if role.Id == 8 {
				// If left end of inner circle
				roleBeforeDemotion = role
				roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
				continue
			} else {
				demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 1)
				if err != nil {
					fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, err)
					return err
				}
				updatedCurrentRoleIds += fmt.Sprintf("%d,", demotedRole.Id)
				roleIdsPostDemote = append(roleIdsPostDemote, demotedRole.Id)
				roleBeforeDemotion = role
				roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
			}
		} else if role.Id > 1 && role.Id < 8 {
			if role.Id == 2 || role.Id == 4 {
				// Server booster role or top contribs - copy across and don't demote from it
				updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
				roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
				roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
			} else {
				// Staff roles
				if role.Id-1 == 2 {
					// Demotion from Moderator leads to being kicked out of the guild
					err = s.GuildMemberDelete(guildId, userId)
					if err != nil {
						fmt.Println("Error kicking member for demoting from Moderator:", err)
						return err
					}
				} else if role.Id-1 == 4 {
					// Demotion from Administrator leads to Moderator
					demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 2)
					if err != nil {
						fmt.Printf("Error ocurred while trying to demote staff role for member with ID %s: %v", userId, err)
						return err
					}
					updatedCurrentRoleIds += fmt.Sprintf("%d,", demotedRole.Id)
					roleIdsPostDemote = append(roleIdsPostDemote, demotedRole.Id)
					roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
				} else {
					demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 1)
					if err != nil {
						fmt.Printf("Error ocurred while trying to demote staff role for member with ID %s: %v", userId, err)
						return err
					}
					updatedCurrentRoleIds += fmt.Sprintf("%d,", demotedRole.Id)
					roleIdsPostDemote = append(roleIdsPostDemote, demotedRole.Id)
					roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
				}
			}
		} else { // Aztec or Arhitect role
			updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
			roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
			roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
		}
	}
	userToUpdate.CurrentRoleIds = updatedCurrentRoleIds

	// Circle and Order (for Inner members)
	currentCircle, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIdsPostDemote)
	userToUpdate.CurrentCircle = currentCircle
	userToUpdate.CurrentInnerOrder = currentOrder

	// Update User in the database
	_, errDemoteUserUpdate := globalsRepo.UsersRepository.UpdateUser(*userToUpdate)
	if errDemoteUserUpdate != nil {
		fmt.Printf("Error ocurred while trying to demote member: %v", errDemoteUserUpdate)
		return err
	}

	// Update Member in the Discord guild
	// Remove all roles
	err = RemoveAllDiscordUserRoles(s, globals.DiscordMainGuildId, userId)
	if err != nil {
		// Revert
		fmt.Printf("An error ocurred while removing all roles for member: %v\n", err)
		err = AddRolesToDiscordUser(s, globals.DiscordMainGuildId, userId, roleIdsPriorDemote)
		if err != nil {
			fmt.Printf("An error ocurred while reverting all roles deletion: %v\n", err)
		}
	}

	// Add new roles
	err = AddRolesToDiscordUser(s, globals.DiscordMainGuildId, userId, roleIdsPostDemote)
	if err != nil {
		fmt.Printf("An error ocurred while adding all roles from DB for member: %v\n", err)
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
			fmt.Println("Error removing role:", err)
			return err
		}
	}

	return nil

}

func AddRolesToDiscordUser(s *discordgo.Session, guildId string, userId string, roleIds []int) error {

	// For each role from the DB
	for _, roleId := range roleIds {
		role, err := globalsRepo.RolesRepository.GetRoleById(roleId)
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
				fmt.Println("Error adding DB role to Discord member:", err)
				return err
			}
		}
	}

	// Process ORDER role post-update (based on the current role state)
	_, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
	fmt.Println("Role IDs", roleIds)
	if currentOrder != nil {
		fmt.Println("Processed Order", *currentOrder)
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
