package member

import (
	"database/sql"
	"fmt"
	"log"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Gets the first staff role of a given member.
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

// Gets the order roles of a given member.
func GetMemberOrderRoles(userId string, defaultOrderRoleNames []string) ([]*dataModels.Role, error) {

	roles, err := globalRepositories.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		return nil, err
	}

	rolesResult := []*dataModels.Role{}
	for _, role := range roles {
		if utils.StringInSlice(role.DisplayName, defaultOrderRoleNames) {
			rolesResult = append(rolesResult, &role)
		}
	}

	return rolesResult, nil

}

// Removes all *order* roles from the database member.
func RemoveAllMemberOrderRoles(userId string, defaultOrderRoleNames []string) error {

	for _, roleName := range defaultOrderRoleNames {
		err := globalRepositories.UsersRepository.RemoveUserRoleWithName(userId, roleName)
		if err != nil {
			return err
		}
	}

	return nil

}

// Removes all roles from the database member.
func RemoveAllMemberRoles(userId string) error {

	err := globalRepositories.UsersRepository.RemoveUserRoles(userId)
	if err != nil {
		return err
	}

	return nil

}

// Returns a string which contains a comma-separated list of role IDs (to be saved in the User entity in the DB),
// an array of integers representing the role IDs as seen in the DB,
// and and error, if applicable.
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
