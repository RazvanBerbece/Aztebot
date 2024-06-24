package member

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func IsVerified(userId string) bool {

	hasAtLeastOneRole := false
	hasVerifiedRole := false
	hasCreatedAtTimestamp := false

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
	}

	roleIds := utils.GetRoleIdsFromRoleString(user.CurrentRoleIds)

	if len(roleIds) > 0 {
		hasAtLeastOneRole = true
	}

	if user.CreatedAt != nil {
		hasCreatedAtTimestamp = true
	}

	for _, roleId := range roleIds {
		if roleId == 1 {
			// role with ID = 1 is always the verified role
			hasVerifiedRole = true
			break
		}
	}

	return (hasVerifiedRole || hasCreatedAtTimestamp) && hasAtLeastOneRole
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

func SetGender(userId string, genderValue string) error {

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		return err
	}

	switch genderValue {
	case "male":
		user.Gender = 0
	case "female":
		user.Gender = 1
	case "nonbin":
		user.Gender = 2
	case "other":
		user.Gender = 3
	default:
		user.Gender = -1
	}

	_, err = globalRepositories.UsersRepository.UpdateUser(*user)
	if err != nil {
		return err
	}

	// Also set gender in leaderboard - if applicable
	count := globalRepositories.MonthlyLeaderboardRepository.EntryExists(userId)
	if count <= 0 {
		if count == -1 {
			return fmt.Errorf("an error ocurred while checking for user leaderboard entry")
		}
	} else {
		err = globalRepositories.MonthlyLeaderboardRepository.UpdateCategoryForUser(userId, user.Gender)
		if err != nil {
			return err
		}
	}

	return nil

}

func GetRep(userId string) (int, error) {

	rep, err := globalRepositories.UserRepRepository.GetRepForUser(userId)
	if err != nil {
		return 0, err
	}

	return rep.Rep, nil

}

// scope: 0 for startup sync; 1 otherwise
func VerifyMember(s *discordgo.Session, logger logging.Logger, guildId string, userId string, scope string) error {

	// Only verify members which haven't been verified yet (according to DB state)
	if !IsVerified(userId) {

		user, err := globalRepositories.UsersRepository.GetUser(userId)
		if err != nil {
			return err
		}

		unixNow := time.Now().Unix()
		user.CreatedAt = &unixNow

		// Newly verified user, so announce in global (if notification channel exists)
		if globalConfiguration.GreetNewVerifiedUsersInChannel {
			if channel, channelExists := globalConfiguration.NotificationChannels["notif-globalGeneralChat"]; channelExists {
				var content string
				if scope == "startup" {
					// at start-up (persistent) sync
					content = fmt.Sprintf("<@%s> has recently verified as an OTA community member! Say hello 🍻", user.UserId)
				} else {
					// at common sync (on role updates, etc.)
					content = fmt.Sprintf("<@%s> has verified as an OTA community member! Say hello 🍻", user.UserId)
				}
				globalMessaging.NotificationsChannel <- events.NotificationEvent{
					TargetChannelId: channel.ChannelId,
					Type:            "DEFAULT",
					TextData:        &content,
				}
			}
		}

		if globalConfiguration.AuditMemberVerificationsInChannel {
			go logger.LogInfo(fmt.Sprintf("`%s` has completed their verification", user.DiscordTag))
		}

		_, updateErr := globalRepositories.UsersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", updateErr)
			return err
		}

		// Give verified role to member
		err = AddDiscordRoleToMember(s, guildId, userId, "Aztec")
		if err != nil {
			log.Println("Error adding  user in DB:", err)
			return err
		}

		return nil

	}

	// Ensure that verified status is visible on Discord too
	err := AddDiscordRoleToMember(s, guildId, userId, "Aztec")
	if err != nil {
		log.Println("Error adding  user in DB:", err)
		return err
	}

	return nil

}
