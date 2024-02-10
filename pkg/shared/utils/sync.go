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

	userExistsResult := globalsRepo.UsersRepository.UserExists(userId)
	switch userExistsResult {
	case -1:
		fmt.Printf("Cannot check whether user %s (%s) exists in the DB during bot startup sync\n", member.User.Username, userId)
	case 0:
		fmt.Printf("Adding new member %s to the OTA DB during bot startup sync\n", member.User.Username)
		var err error
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
	case 1:
		// Already exists
		var err error
		user, err = globalsRepo.UsersRepository.GetUser(userId)
		if err != nil {
			log.Fatalf("Error ocurred retrieving user from the DB: %v\n", err)
			return err
		}
	}

	// At this point, user already exists because it was inserted at sync time,
	// or it exists because it was previously synced
	if user != nil {

		// Sync all other user details between the Discord server and the database (mostly updating the DB with Discord data)
		// Get current roles from user (as they appear on the Discord guild)
		currentRoleIds, roleIds, err := GetUserRolesFromDiscord(s, guildId, *user, *member)
		if err != nil {
			log.Println("Error retrieving user's roles as DB data from the Discord Guild:", err)
			return err
		}

		// `Aztec` verification -- user has Aztec role and is verified
		for _, roleId := range roleIds {
			if roleId == 1 && user.CreatedAt == nil {
				unixNow := time.Now().Unix()
				user.CreatedAt = &unixNow
			}
		}

		user.CurrentRoleIds = currentRoleIds

		// Circle and Order (for Inner members)
		currentCircle, currentOrder := GetCircleAndOrderForGivenRoles(roleIds)
		user.CurrentCircle = currentCircle
		user.CurrentInnerOrder = currentOrder

		_, updateErr := usersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", updateErr)
			return err
		}

		return nil
	}

	return fmt.Errorf("no update was executed")
}

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild. This function uses the shared global DB connections.
func SyncUser(s *discordgo.Session, guildId string, userId string, member *discordgo.Member) error {

	var user *dataModels.User

	userExistsResult := globalsRepo.UsersRepository.UserExists(userId)
	switch userExistsResult {
	case -1:
		fmt.Printf("Cannot check whether user %s (%s) exists in the DB\n", member.User.Username, userId)
	case 0:
		fmt.Printf("Adding new member %s to the OTA DB\n", member.User.Username)
		var err error
		user, err = globalsRepo.UsersRepository.SaveInitialUserDetails(member.User.Username, userId)
		if err != nil {
			log.Fatalf("Cannot store initial user %s with id %s: %v\n", member.User.Username, userId, err)
			return err
		}
		errStats := globalsRepo.UserStatsRepository.SaveInitialUserStats(userId)
		if errStats != nil {
			log.Fatalf("Cannot store initial user stats: %v\n", errStats)
			return errStats
		}
	case 1:
		// Already exists
		var err error
		user, err = globalsRepo.UsersRepository.GetUser(userId)
		if err != nil {
			log.Fatalf("Error ocurred retrieving user from the DB: %v\n", err)
			return err
		}
	}

	// At this point, user already exists because it was inserted at sync time,
	// or it exists because it was previously synced
	if user != nil {

		// UPDATE Roles and Role After-Effects
		// Get current roles for user (as they appear on the Discord guild and found in the Roles DB)
		currentRoleIds, roleIds, err := GetUserRolesFromDiscord(s, guildId, *user, *member)
		if err != nil {
			log.Println("Error retrieving user's roles as DB data from the Discord Guild:", err)
			return err
		}

		// `Aztec` verification -- user has Aztec role and is verified
		for _, roleId := range roleIds {
			if roleId == 1 && user.CreatedAt == nil {
				unixNow := time.Now().Unix()
				user.CreatedAt = &unixNow
			}
		}

		user.CurrentRoleIds = currentRoleIds

		// Circle and Order (for Inner members)
		currentCircle, currentOrder := GetCircleAndOrderForGivenRoles(roleIds)
		user.CurrentCircle = currentCircle
		user.CurrentInnerOrder = currentOrder

		// Save changes
		_, updateErr := globalsRepo.UsersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", updateErr)
			return err
		}

		return nil
	}

	return fmt.Errorf("no update was executed")
}

// Returns a string which contains a comma-separated list of role IDs (to be saved in the User entity in the DB),
// an array of integers representing the role IDs as seen in the DB,
// and error, if applicable.
func GetUserRolesFromDiscord(s *discordgo.Session, guildId string, user dataModels.User, member discordgo.Member) (string, []int, error) {

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

		roleDax, err := globalsRepo.RolesRepository.GetRole(userRoleObj.Name)
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
