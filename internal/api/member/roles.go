package member

import (
	"fmt"
	"log"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Removes all roles from the database member.
func RemoveAllMemberRoles(userId string) error {

	err := globalsRepo.UsersRepository.RemoveUserRoles(userId)
	if err != nil {
		return err
	}

	return nil

}

func IsStaff(userId string, staffRoles []string) bool {

	roles, err := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
	}

	for _, role := range roles {
		if utils.StringInSlice(role.DisplayName, staffRoles) {
			return true
		}
	}

	return false
}

// Allows demotion either from a staff role or an inner order role.
func DemoteMember(s *discordgo.Session, guildId string, userId string, demoteType string) error {

	userToUpdate, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, err)
		return err
	}

	// DEMOTE STRATEGY (BOTH INNER CIRCLE ORDERS AND STAFF ROLE DEMOTIONS)
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
			if demoteType == "STAFF" {
				updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
				roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
				roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
			} else {
				if role.Id == 8 {
					// If left end of inner circle
					roleBeforeDemotion = role
					roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
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
			}
		} else if role.Id > 1 && role.Id < 8 {
			if role.Id == 2 || role.Id == 4 {
				// Server booster role or top contribs - copy across and don't demote from it
				updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
				roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
				roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
			} else {
				// Staff roles
				if demoteType == "ORDER" {
					updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
					roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
					roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
				} else {
					if role.Id-1 == 2 {
						// Demotion from Moderator leads to being kicked out of the guild
						err = KickMember(s, guildId, userId)
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
