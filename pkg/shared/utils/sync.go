package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild. It uses repositories injected via the argument list to prevent connection attempt floods.
func SyncUserPersistent(s *discordgo.Session, guildId string, userId string, member *discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) error {

	var user *dataModels.User

	userExistsResult, err := globalsRepo.UsersRepository.UserExists(userId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Cannot retrieve user %s with id %s during bot startup sync time: %v", member.User.Username, userId, err)
		return err
	} else if err != nil && err == sql.ErrNoRows && !userExistsResult {
		fmt.Printf("Adding new member %s to the OTA DB during bot startup sync\n", member.User.Username)
		user, err = globalsRepo.UsersRepository.SaveInitialUserDetails(member.User.Username, userId)
		if err != nil {
			log.Fatalf("Cannot store initial user %s with id %s during bot startup sync: %v", member.User.Username, userId, err)
			return err
		}
		errStats := userStatsRepository.SaveInitialUserStats(userId)
		if errStats != nil {
			log.Fatalf("Cannot store initial user stats during bot startup sync: %v", err)
			return errStats
		}
	}

	if user != nil {

		// Setup user stats if the user doesn't have an entity in UserStats
		userExistsResult, err := globalsRepo.UsersRepository.UserExists(userId)
		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("Cannot retrieve user stats %s with id %s during bot startup sync time: %v", member.User.Username, userId, err)
			return err
		} else if err != nil && err == sql.ErrNoRows && !userExistsResult {
			errStatsInit := userStatsRepository.SaveInitialUserStats(userId)
			if errStatsInit != nil {
				log.Fatalf("Cannot store initial user %s stats during bot startup sync time: %v", member.User.Username, errStatsInit)
				return errStatsInit
			}
		}

		// Sync all other user details between the Discord server and the database (mostly updating the DB with Discord data)
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

		_, updateErr := usersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", err)
			return err
		}

		// fmt.Printf("Synced user %s\n", updatedUser.DiscordTag)

		return nil
	}

	return fmt.Errorf("no update was executed")
}

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild. This function uses the shared global DB connections.
func SyncUser(s *discordgo.Session, guildId string, userId string, member *discordgo.Member) error {

	var user *dataModels.User

	userExistsResult, err := globalsRepo.UsersRepository.UserExists(userId)
	if err != nil && err == sql.ErrNoRows && !userExistsResult {
		fmt.Printf("Adding new member %s to the OTA DB\n", member.User.Username)
		user, err = globalsRepo.UsersRepository.SaveInitialUserDetails(member.User.Username, userId)
		if err != nil {
			log.Fatalf("Cannot store user %s with id %s: %v", member.User.Username, userId, err)
			return err
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
			roleDax, err := globalsRepo.RolesRepository.GetRole(userRoleObj.Name)
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

		_, updateErr := globalsRepo.UsersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", err)
			return err
		}

		return nil
	}

	return fmt.Errorf("no update was executed")
}
