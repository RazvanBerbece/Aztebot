package member

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func IsFullyVerified(userId string) bool {

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

func ProcessMemberVerification(s *discordgo.Session, logger logging.Logger, usersRepository repositories.UsersRepository, jailRepository repositories.JailRepository, guildId string, userId string, currentRoleIds []int, scope string) error {

	// Skip this if member is jailed
	jailed, _ := jailRepository.GetJailedUser(userId)
	if jailed != nil {
		return nil
	}

	user, err := usersRepository.GetUser(userId)
	if err != nil {
		return err
	}

	if user.CreatedAt == nil { // No member join timestamp
		if !utils.IntInSlice(globalConfiguration.DefaultVerifiedRoleId, currentRoleIds) {
			// ignore
			return nil
		} else { // Member obtained the verified role but not the timestamp

			unixNow := time.Now().Unix()
			err := usersRepository.SetUserCreatedAt(userId, unixNow)
			if err != nil {
				log.Println("Error updating user in DB:", err)
				return err
			}

			// Global channel announcement
			if globalConfiguration.GreetNewVerifiedUsersInChannel {
				if channel, channelExists := globalConfiguration.NotificationChannels["notif-globalGeneralChat"]; channelExists {
					var content string
					if scope == "startup" {
						// at start-up (persistent) sync
						content = fmt.Sprintf("<@%s> has recently been verified as an OTA community member! Say hello 🍻", user.UserId)
					} else if scope == "default" {
						// at common sync (on role updates, etc.)
						content = fmt.Sprintf("<@%s> has just been verified as an OTA community member! Say hello 🍻", user.UserId)
					}
					globalMessaging.NotificationsChannel <- events.NotificationEvent{
						TargetChannelId: channel.ChannelId,
						Type:            "DEFAULT",
						TextData:        &content,
					}
				}
			}

			// Auditing
			if globalConfiguration.AuditMemberVerificationsInChannel {
				if scope == "startup" {
					go logger.LogInfo(fmt.Sprintf("`%s` has completed their verification during STARTUP", user.DiscordTag))
				} else {
					go logger.LogInfo(fmt.Sprintf("`%s` has completed their verification during BAU", user.DiscordTag))
				}
			}
		}
	} else { // Existing member join timestamp
		if !utils.IntInSlice(globalConfiguration.DefaultVerifiedRoleId, currentRoleIds) {
			// Member has timestamp but not the actual role, consider as intentional and skip
			return nil
		}
	}

	return nil

}
