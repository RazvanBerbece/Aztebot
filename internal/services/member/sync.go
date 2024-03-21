package member

import (
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild. This function uses the shared global DB connections.
func SyncMember(s *discordgo.Session, guildId string, userId string, member *discordgo.Member) error {

	var user *dataModels.User

	userExistsResult := globalRepositories.UsersRepository.UserExists(userId)
	switch userExistsResult {
	case -1:
		fmt.Printf("Cannot check whether user %s (%s) exists in the DB\n", member.User.Username, userId)
	case 0:
		var err error
		user, err = globalRepositories.UsersRepository.SaveInitialUserDetails(member.User.Username, userId)
		if err != nil {
			log.Fatalf("Cannot store initial user %s with id %s: %v\n", member.User.Username, userId, err)
			return err
		}
		errStats := globalRepositories.UserStatsRepository.SaveInitialUserStats(userId)
		if errStats != nil {
			log.Fatalf("Cannot store initial user stats: %v\n", errStats)
			return errStats
		}
		fmt.Printf("Added new member entries %s to the OTA DB\n", member.User.Username)
	case 1:
		// Already exists
		var err error
		user, err = globalRepositories.UsersRepository.GetUser(userId)
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
		currentRoleIds, roleIds, err := GetMemberRolesFromDiscordAsLocalIdList(s, guildId, *user, *member)
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

		// Circle and Order (for Inner members) in the DB
		currentCircle, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
		user.CurrentCircle = currentCircle
		user.CurrentInnerOrder = currentOrder

		// Save changes
		_, updateErr := globalRepositories.UsersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", updateErr)
			return err
		}

		return nil
	}

	return fmt.Errorf("no update was executed")
}

// Takes in a discord member and syncs the database User with the current member details
// as they appear on the Discord guild. It uses repositories injected via the argument list to prevent connection attempt floods.
func SyncMemberPersistent(s *discordgo.Session, guildId string, userId string, member *discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) error {

	var user *dataModels.User

	userExistsResult := globalRepositories.UsersRepository.UserExists(userId)
	switch userExistsResult {
	case -1:
		fmt.Printf("Cannot check whether user %s (%s) exists in the DB during bot startup sync\n", member.User.Username, userId)
	case 0:
		var err error
		user, err = globalRepositories.UsersRepository.SaveInitialUserDetails(member.User.Username, userId)
		if err != nil {
			log.Fatalf("Cannot store initial user %s with id %s during bot startup sync: %v", member.User.Username, userId, err)
			return err
		}
		errStats := userStatsRepository.SaveInitialUserStats(userId)
		if errStats != nil {
			log.Fatalf("Cannot store initial user stats during bot startup sync: %v", err)
			return errStats
		}
		fmt.Printf("Added new member entries %s to the OTA DB during bot startup sync\n", member.User.Username)
	case 1:
		// Already exists
		var err error
		user, err = globalRepositories.UsersRepository.GetUser(userId)
		if err != nil {
			log.Fatalf("Error ocurred retrieving user from the DB: %v\n", err)
			return err
		}
		// Check whether user has user stats entity
		userStatsExistsResult := globalRepositories.UserStatsRepository.UserStatsExist(userId)
		switch userStatsExistsResult {
		case -1:
			// Error ocurred
			fmt.Printf("Cannot check whether user %s (%s) exists in the DB during bot startup sync\n", member.User.Username, userId)
		case 0:
			// Stats don't exist
			err = globalRepositories.UserStatsRepository.SaveInitialUserStats(userId)
			if err != nil {
				log.Printf("Failed to store initial user stats at startup: %v", err)
				return err
			}
		case 1:
			// Stats exist
			// pass
		}
	}

	// At this point, user already exists because it was inserted at sync time,
	// or it exists because it was previously synced
	if user != nil {

		// Sync all other user details between the Discord server and the database (mostly updating the DB with Discord data)
		// Get current roles from user (as they appear on the Discord guild)
		currentRoleIds, roleIds, err := GetMemberRolesFromDiscordAsLocalIdList(s, guildId, *user, *member)
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
		currentCircle, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
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
