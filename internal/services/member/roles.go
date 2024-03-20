package member

import (
	"database/sql"
	"fmt"
	"log"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func GetMemberStaffRole(userId string, staffRoles []string) (*dataModels.Role, error) {

	roles, err := globalRepositories.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		return nil, err
	}

	for _, role := range roles {
		if utils.StringInSlice(role.DisplayName, staffRoles) {
			return &role, nil
		}
	}

	return nil, nil

}

// Removes all roles from the database member.
func RemoveAllMemberRoles(userId string) error {

	err := globalRepositories.UsersRepository.RemoveUserRoles(userId)
	if err != nil {
		return err
	}

	return nil

}

func IsStaff(userId string, staffRoles []string) bool {

	roles, err := globalRepositories.UsersRepository.GetRolesForUser(userId)
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

	userToUpdate, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, err)
		return err
	}

	// DEMOTE STRATEGY (BOTH INNER CIRCLE ORDERS AND STAFF ROLE DEMOTIONS)
	userRoles, errUsrRole := globalRepositories.UsersRepository.GetRolesForUser(userId)
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
					demotedRole, err := globalRepositories.RolesRepository.GetRoleById(role.Id - 1)
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
						demotedRole, err := globalRepositories.RolesRepository.GetRoleById(role.Id - 2)
						if err != nil {
							fmt.Printf("Error ocurred while trying to demote staff role for member with ID %s: %v", userId, err)
							return err
						}
						updatedCurrentRoleIds += fmt.Sprintf("%d,", demotedRole.Id)
						roleIdsPostDemote = append(roleIdsPostDemote, demotedRole.Id)
						roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
					} else {
						demotedRole, err := globalRepositories.RolesRepository.GetRoleById(role.Id - 1)
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
	_, errDemoteUserUpdate := globalRepositories.UsersRepository.UpdateUser(*userToUpdate)
	if errDemoteUserUpdate != nil {
		fmt.Printf("Error ocurred while trying to demote member: %v", errDemoteUserUpdate)
		return err
	}

	// Update Member in the Discord guild
	// Remove all roles
	err = RemoveAllDiscordRolesFromMember(s, globalConfiguration.DiscordMainGuildId, userId)
	if err != nil {
		// Revert
		fmt.Printf("An error ocurred while removing all roles for member: %v\n", err)
		err = AddDiscordRolesToMember(s, globalConfiguration.DiscordMainGuildId, userId, roleIdsPriorDemote)
		if err != nil {
			fmt.Printf("An error ocurred while reverting all roles deletion: %v\n", err)
		}
	}

	// Add new roles
	err = AddDiscordRolesToMember(s, globalConfiguration.DiscordMainGuildId, userId, roleIdsPostDemote)
	if err != nil {
		fmt.Printf("An error ocurred while adding all roles from DB for member: %v\n", err)
	}

	return nil

}

// Returns a string which contains a comma-separated list of role IDs (to be saved in the User entity in the DB),
// an array of integers representing the role IDs as seen in the DB,
// and error, if applicable.
func GetMemberRolesFromDiscordAsLocalIdList(s *discordgo.Session, guildId string, user dataModels.User, member discordgo.Member) (string, []int, error) {

	var currentRoleIds string // string representing a list of role IDs (this is to be stored in the DB)
	var roleIds []int         // integer list of the role IDs (like above, but an array of int IDs)

	// Build a list of roles taken from the Discord guild
	// and then use the list to update the role IDs, circle and order in the database for the given user & member pair
	for _, role := range member.Roles {

		userRoleObj, err := s.State.Role(guildId, role) // role DisplayName in OTA DB
		if err != nil {
			log.Println("Error getting role from Discord servers:", err)
			return "", nil, err
		}

		roleDax, err := globalRepositories.RolesRepository.GetRole(userRoleObj.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				// This will probably be a role which is assigned to the three orders or something, so we can ignore
				// and move on to the other roles of the user
				continue
			} else {
				return "", nil, err
			}
		} else {
			// Build up the role data for the current member as observed in the Roles DB
			currentRoleIds += fmt.Sprintf("%d,", roleDax.Id)
			roleIds = append(roleIds, roleDax.Id)
		}
	}

	return currentRoleIds, roleIds, nil

}
