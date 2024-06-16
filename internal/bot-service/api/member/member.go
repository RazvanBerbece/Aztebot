package member

import (
	"fmt"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
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

	// DEMOTE STRATEGY (ATM CONSIDER THAT MEMBER IS PART OF STAFF SO DEMOTE WORKS IN THE INNER CIRCLE ORDERS)
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
	var rolePostDemotion *dataModels.Role
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
				rolePostDemotion = demotedRole
				roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
			}
		} else {
			updatedCurrentRoleIds += fmt.Sprintf("%d,", role.Id)
			roleIdsPostDemote = append(roleIdsPostDemote, role.Id)
			roleIdsPriorDemote = append(roleIdsPriorDemote, role.Id)
		}
	}
	userToUpdate.CurrentRoleIds = updatedCurrentRoleIds

	// Circle and Order (for Inner members)
	_, previousOrder := utils.GetCircleAndOrderForGivenRoles(roleIdsPriorDemote)
	currentCircle, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIdsPostDemote)
	userToUpdate.CurrentCircle = currentCircle
	userToUpdate.CurrentInnerOrder = currentOrder

	_, errDemoteUserUpdate := globalsRepo.UsersRepository.UpdateUser(*userToUpdate)
	if errDemoteUserUpdate != nil {
		fmt.Printf("Error ocurred while trying to demote member: %v", errDemoteUserUpdate)
		return err
	}

	// DEMOTE ON THE DISCORD SERVER (UPDATE ACTUAL ROLES: ROLE, CIRCLE, ORDER)
	oldDiscordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleBeforeDemotion.DisplayName)
	if oldDiscordRoleId != nil {
		// Remove old role (previous to the demotion)
		err = s.GuildMemberRoleRemove(guildId, userId, *oldDiscordRoleId)
		if err != nil {
			fmt.Println("Error removing role:", err)
			return err
		}
		// Add new role (post demotion) if necessary
		if rolePostDemotion != nil {
			newDiscordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, rolePostDemotion.DisplayName)
			err = s.GuildMemberRoleAdd(guildId, userId, *newDiscordRoleId)
			if err != nil {
				fmt.Println("Error adding role:", err)
				return err
			}
		}
		// Remove old order role if necessary
		if currentOrder != previousOrder {
			var prevOrderValue int = *previousOrder
			var discordOrderRoleIdToDelete *string
			if prevOrderValue == 3 {
				discordOrderRoleIdToDelete = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
			} else if prevOrderValue == 2 {
				discordOrderRoleIdToDelete = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
			} else if prevOrderValue == 1 {
				discordOrderRoleIdToDelete = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
			}
			err = s.GuildMemberRoleRemove(guildId, userId, *discordOrderRoleIdToDelete)
			if err != nil {
				fmt.Println("Error removing role:", err)
				return err
			}
			// Add new order role if necessary
			if currentOrder != nil {
				var discordOrderRoleIdToAdd *string
				if *currentOrder == 3 {
					discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
				} else if *currentOrder == 2 {
					discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
				} else if *currentOrder == 1 {
					discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
				}
				err = s.GuildMemberRoleAdd(guildId, userId, *discordOrderRoleIdToAdd)
				if err != nil {
					fmt.Println("Error removing role:", err)
					return err
				}
			}
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
	var roleID string
	for _, role := range guild.Roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	return &roleID
}
