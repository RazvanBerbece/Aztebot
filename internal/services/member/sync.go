package member

import (
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/metrics"
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
		user, err = globalRepositories.UsersRepository.SaveInitialUserDetails(member.User.Username, userId, nil)
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
				// Newly verified user, so announce in global (if notification channel exists)
				if channel, channelExists := globalConfiguration.NotificationChannels["notif-globalGeneralChat"]; channelExists {
					content := fmt.Sprintf("<@%s> has joined the OTA community! Say hello 🍻", user.UserId)
					globalMessaging.NotificationsChannel <- events.NotificationEvent{
						TargetChannelId: channel.ChannelId,
						Type:            "DEFAULT",
						TextData:        &content,
					}
				}
				break
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
func SyncMemberPersistent(s *discordgo.Session, guildId string, userId string, member *discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository, defaultOrderRoleNames []string, syncProgression bool) error {

	// TODO: Remove this line
	defer metrics.Duration(metrics.Track("SyncMemberPersistent"))

	var user *dataModels.User
	var userStats *dataModels.UserStats

	userExistsResult := globalRepositories.UsersRepository.UserExists(userId)
	switch userExistsResult {
	case -1:
		fmt.Printf("Cannot check whether user %s (%s) exists in the DB during bot startup sync\n", member.User.Username, userId)
	case 0:
		var err error
		user, err = globalRepositories.UsersRepository.SaveInitialUserDetails(member.User.Username, userId, nil)
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
			userStats, err = globalRepositories.UserStatsRepository.GetStatsForUser(userId)
			if err != nil {
				log.Fatalf("Error ocurred retrieving user stats from the DB: %v\n", err)
				return err
			}
		}
	}

	// At this point, user and stats already exist because they were inserted at sync time,
	// or they exists becaus they were previously synced
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
				// Newly verified user, so announce in global (if notification channel exists)
				if channel, channelExists := globalConfiguration.NotificationChannels["notif-globalGeneralChat"]; channelExists {
					content := fmt.Sprintf("<@%s> has recently joined the OTA community! Say hello 🍻", user.UserId)
					globalMessaging.NotificationsChannel <- events.NotificationEvent{
						TargetChannelId: channel.ChannelId,
						Type:            "DEFAULT",
						TextData:        &content,
					}
				}
				break
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

		if syncProgression {
			err = ResolveProgressionMismatchForMember(s, guildId, userId, user.CurrentExperience, userStats.NumberMessagesSent, userStats.TimeSpentInVoiceChannels, defaultOrderRoleNames)
			if err != nil {
				log.Println("Error syncing progression for member:", updateErr)
				return err
			}
			err = RefreshDiscordOrderRoleForMember(s, guildId, userId)
			if err != nil {
				log.Println("Error refreshing order role for member:", updateErr)
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("no update was executed")
}

func ResolveProgressionMismatchForMember(s *discordgo.Session, userGuildId string, userId string, userXp float64, userNumberMessagesSent int, userTimeSpentInVc int, defaultOrderRoleNames []string) error {

	// don't sync progression for unverified users
	if !IsVerified(userId) {
		return nil
	}

	// Check current stats against progression table
	// Figure out the promoted role to be given
	processedRoleName, processedLevel := GetRoleNameAndLevelFromStats(userXp, userNumberMessagesSent, userTimeSpentInVc)

	currentOrderRoles, err := GetMemberOrderRoles(userId, defaultOrderRoleNames)
	if err != nil {
		fmt.Printf("Error occurred while reading member order role from DB: %v\n", err)
		return err
	}

	// Solve mismatches where the member has a rank on the server but shouldn't
	// according to the progression rules (type 1, 2, 3, 4)
	if processedLevel == 0 && processedRoleName == "" && len(currentOrderRoles) > 0 {
		// mismatch, need to reset
		err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
		if err != nil {
			fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
			return err
		}

		for _, orderRole := range currentOrderRoles {
			err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, orderRole.Id)
			if err != nil {
				fmt.Printf("Error occurred while removing member role from DB: %v\n", err)
			}
		}

		user, err := globalRepositories.UsersRepository.GetUser(userId)
		if err != nil {
			fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
		}
		err = RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
		if err != nil {
			fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
		}

		fmt.Printf("Mismatch (type 1) for %s resolved.\n", user.DiscordTag)
	} else if processedLevel > 0 && processedRoleName != "" && len(currentOrderRoles) == 1 {
		if currentOrderRoles[0].DisplayName != processedRoleName {
			// Solve mismatches where the member has a rank on the server but their
			// actual non-zero rank is different (type 2)
			err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
			if err != nil {
				fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
				return err
			}

			err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, currentOrderRoles[0].Id)
			if err != nil {
				fmt.Printf("Error occurred while removing member role from DB: %v\n", err)
			}

			promotedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
			if err != nil {
				fmt.Printf("Error occurred while reading role from DB: %v\n", err)
				return err
			}

			err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
			if err != nil {
				fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
			}

			user, err := globalRepositories.UsersRepository.GetUser(userId)
			if err != nil {
				fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
			}
			err = RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
			if err != nil {
				fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
			}

			fmt.Printf("Mismatch (type 2) for %s resolved.\n", user.DiscordTag)
		}
	} else if processedLevel > 0 && processedRoleName != "" && len(currentOrderRoles) > 1 {
		// Solve mismatches where the member has multiple ranks on the server but their
		// actual non-zero rank is different (type 3)
		for _, role := range currentOrderRoles {
			if role.DisplayName != processedRoleName {
				err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, role.Id)
				if err != nil {
					fmt.Printf("Error occurred while removing member role from DB: %v\n", err)
				}
			}
		}

		err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
		if err != nil {
			fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
			return err
		}

		promotedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
		if err != nil {
			fmt.Printf("Error occurred while reading role from DB: %v\n", err)
			return err
		}

		err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
		if err != nil {
			fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
		}

		user, err := globalRepositories.UsersRepository.GetUser(userId)
		if err != nil {
			fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
		}
		err = RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
		if err != nil {
			fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
		}

		fmt.Printf("Mismatch (type 3) for %s resolved.\n", user.DiscordTag)
	} else if processedLevel > 0 && processedRoleName != "" && len(currentOrderRoles) == 0 {
		// Solve mismatches where the member has no rank on the server but their
		// actual rank is different and non-zero (type 4)
		err := globalRepositories.UsersRepository.SetLevel(userId, processedLevel)
		if err != nil {
			fmt.Printf("Error occurred while setting member level in DB: %v\n", err)
			return err
		}

		promotedRole, err := globalRepositories.RolesRepository.GetRole(processedRoleName) // to append
		if err != nil {
			fmt.Printf("Error occurred while reading role from DB: %v\n", err)
			return err
		}

		err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
		if err != nil {
			fmt.Printf("Error occurred while appending role ID to member in DB: %v\n", err)
		}

		user, err := globalRepositories.UsersRepository.GetUser(userId)
		if err != nil {
			fmt.Printf("Error occurred while retrieving user and roles from DB: %v\n", err)
		}
		err = RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
		if err != nil {
			fmt.Printf("Error occurred while refreshing member roles on-Discord: %v\n", err)
		}

		fmt.Printf("Mismatch (type 4) for %s resolved.\n", user.DiscordTag)
	}

	return nil
}
