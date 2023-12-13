package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/bwmarrin/discordgo"
)

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild.
func SyncUser(s *discordgo.Session, guildId string, userId string, member *discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository) error {

	user, err := usersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
		if err == sql.ErrNoRows {
			log.Printf("Storing user with id %s", userId)
			user, err = usersRepository.SaveInitialUserDetails(member.User.Username, userId)
			if err != nil {
				log.Fatalf("Cannot store user %s with id %s: %v", member.User.Username, userId, err)
				return err
			}
		}
	}

	if user != nil {
		// Roles
		// Get current roles from user (as they appear on the Discord guild)
		var currentRoleIds string
		var roleIds []int
		for _, role := range member.Roles {
			// Build a list of roles taken from the Discord guild
			// and then use the list to update the role IDs, circle and order in the database
			userRoleObj, err := s.State.Role(guildId, role) // role DisplayName in OTA DB
			if err != nil {
				log.Println("Error getting role from Discord servers:", err)
				return err
			}
			roleDax, err := rolesRepository.GetRole(userRoleObj.Name)
			if err != nil {
				if err == sql.ErrNoRows {
					// This will probably be a role which is assigned to the three orders or something, so we can ignore
					// and move on to the other roles of the user
					continue
				} else {
					return err
				}
			} else {
				// `Aztec` verification
				if roleDax.Id == 1 && user.CreatedAt == nil {
					unixNow := time.Now().Unix()
					user.CreatedAt = &unixNow
				}
				// Role IDs
				currentRoleIds += fmt.Sprintf("%d,", roleDax.Id)
				// Circle
				roleIds = append(roleIds, roleDax.Id)
			}
		}

		user.CurrentRoleIds = currentRoleIds

		var hasInnerCircleId bool = false
		var maxInnerOrderId int = -1
		for _, roleId := range roleIds {
			circle, order := GetCircleAndOrderFromRoleId(roleId)
			if circle == 1 {
				hasInnerCircleId = true
				if order > maxInnerOrderId {
					maxInnerOrderId = order
				}
			}
		}

		if hasInnerCircleId {
			user.CurrentCircle = "INNER"
		} else {
			user.CurrentCircle = "OUTER"
		}

		if maxInnerOrderId == -1 {
			user.CurrentInnerOrder = nil
		} else {
			user.CurrentInnerOrder = &maxInnerOrderId
		}

		updatedUser, updateErr := usersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", err)
			return err
		}
		fmt.Printf("Synced user %s\n", updatedUser.DiscordTag)

		return nil
	}

	return fmt.Errorf("no update was executed")
}
