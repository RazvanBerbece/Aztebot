package member

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/api/notifications"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func IsStaff(userId string) bool {

	roles, err := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
	}

	for _, role := range roles {
		if utils.RoleIsStaffRole(role.Id) {
			return true
		}
	}

	return false
}

func KickMember(s *discordgo.Session, guildId string, userId string) error {
	// Delete member from server
	err := s.GuildMemberDelete(guildId, userId)
	if err != nil {
		fmt.Println("Error kicking member from guild:", err)
		return err
	}
	// Delete member-related entries from the databases
	err = DeleteAllMemberData(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s data from DB tables: %v", userId, err)
		return err
	}
	return nil
}

func DeleteAllMemberData(userId string) error {
	err := globalsRepo.UserStatsRepository.DeleteUserStats(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s stats from DB: %v", userId, err)
	}
	err = globalsRepo.UsersRepository.DeleteUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s from DB: %v", userId, err)
	}
	err = globalsRepo.WarnsRepository.DeleteAllWarningsForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s warnings from DB: %v", userId, err)
	}
	err = globalsRepo.TimeoutsRepository.ClearTimeoutForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s active timeouts from DB: %v", userId, err)
	}
	err = globalsRepo.TimeoutsRepository.ClearArchivedTimeoutsForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s archived timeouts from DB: %v", userId, err)
	}
	return nil
}

// Allows demotion either from a staff role or an inner order role.
func DemoteMember(s *discordgo.Session, guildId string, userId string, demoteType string) error {

	userToUpdate, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Error ocurred while trying to demote member with ID %s: %v", userId, err)
		return err
	}

	// DEMOTE STRATEGY (BOTH INNER CIRCLE ORDERS AND STAFF ROLE DEMOTIONS)
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
					demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 1)
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
						demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 2)
						if err != nil {
							fmt.Printf("Error ocurred while trying to demote staff role for member with ID %s: %v", userId, err)
							return err
						}
						updatedCurrentRoleIds += fmt.Sprintf("%d,", demotedRole.Id)
						roleIdsPostDemote = append(roleIdsPostDemote, demotedRole.Id)
						roleIdsPriorDemote = append(roleIdsPriorDemote, roleBeforeDemotion.Id)
					} else {
						demotedRole, err := globalsRepo.RolesRepository.GetRoleById(role.Id - 1)
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
	_, errDemoteUserUpdate := globalsRepo.UsersRepository.UpdateUser(*userToUpdate)
	if errDemoteUserUpdate != nil {
		fmt.Printf("Error ocurred while trying to demote member: %v", errDemoteUserUpdate)
		return err
	}

	// Update Member in the Discord guild
	// Remove all roles
	err = RemoveAllDiscordUserRoles(s, globals.DiscordMainGuildId, userId)
	if err != nil {
		// Revert
		fmt.Printf("An error ocurred while removing all roles for member: %v\n", err)
		err = AddRolesToDiscordUser(s, globals.DiscordMainGuildId, userId, roleIdsPriorDemote)
		if err != nil {
			fmt.Printf("An error ocurred while reverting all roles deletion: %v\n", err)
		}
	}

	// Add new roles
	err = AddRolesToDiscordUser(s, globals.DiscordMainGuildId, userId, roleIdsPostDemote)
	if err != nil {
		fmt.Printf("An error ocurred while adding all roles from DB for member: %v\n", err)
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
	var roleID string = ""
	for _, role := range guild.Roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		fmt.Println("No role ID was found for role name", roleName)
		return nil
	}

	return &roleID
}

func RemoveAllDiscordUserRoles(s *discordgo.Session, guildId string, userId string) error {

	// Get the member's roles
	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return err
	}

	// Find all user's roles and delete them
	for _, roleID := range member.Roles {
		err = s.GuildMemberRoleRemove(guildId, userId, roleID)
		if err != nil {
			fmt.Println("Error removing role:", err)
			return err
		}
	}

	return nil

}

func GiveDiscordRoleToMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Add the role by role ID to the Discord member
	err := s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Println("Error adding giving to Discord member:", err)
		return err
	}

	return nil

}

func RemoveDiscordRoleFromMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Remove the role by role ID from the Discord member
	err := s.GuildMemberRoleRemove(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Println("Error adding giving to Discord member:", err)
		return err
	}

	return nil

}

func AddRolesToDiscordUser(s *discordgo.Session, guildId string, userId string, roleIds []int) error {

	// For each role from the DB
	for _, roleId := range roleIds {
		role, err := globalsRepo.RolesRepository.GetRoleById(roleId)
		if err != nil {
			fmt.Printf("Error ocurred while adding DB roles to Discord member: %v\n", err)
			return err
		}
		// Get the role ID by display name from Discord
		discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, role.DisplayName)
		if discordRoleId != nil {
			// Add the role by role ID to the Discord member
			err = s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
			if err != nil {
				fmt.Println("Error adding DB role to Discord member:", err)
				return err
			}
		}
	}

	// Process ORDER role post-update (based on the current role state)
	_, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
	if currentOrder != nil {
		var discordOrderRoleIdToAdd *string
		if *currentOrder == 3 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
		} else if *currentOrder == 2 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
		} else if *currentOrder == 1 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
		}
		err := s.GuildMemberRoleAdd(guildId, userId, *discordOrderRoleIdToAdd)
		if err != nil {
			fmt.Println("Error adding order role to member:", err)
			return err
		}
	}

	return nil

}

func GetMemberXpRank(userId string) (*int, error) {

	xpRank, err := globalsRepo.UserStatsRepository.GetUserXpRank(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard XP rank for user %s", userId)
		return nil, err
	}

	return xpRank, nil
}

func GetMemberRankInLeaderboards(userId string) (map[string]int, error) {

	results := make(map[string]int)

	// Get place in the messages sent leaderboard
	msgRank, err := globalsRepo.UserStatsRepository.GetUserLeaderboardRank(userId, "msg")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard msg rank for user %s", userId)
		return nil, err
	}
	if *msgRank != 0 {
		results["msg"] = *msgRank
	}
	// Get place in the reactions received leaderboard
	reactRank, err := globalsRepo.UserStatsRepository.GetUserLeaderboardRank(userId, "react")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard react rank for user %s", userId)
		return nil, err
	}
	if *reactRank != 0 {
		results["react"] = *reactRank
	}
	// Get place in the time spent in VCs leaderboard
	vcRank, err := globalsRepo.UserStatsRepository.GetUserLeaderboardRank(userId, "vc")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard vc rank for user %s", userId)
		return nil, err
	}
	if *vcRank != 0 {
		results["vc"] = *vcRank
	}
	// Get place in the time spent in music channels leaderboard
	musicRank, err := globalsRepo.UserStatsRepository.GetUserLeaderboardRank(userId, "music")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard music rank for user %s", userId)
		return nil, err
	}
	if *musicRank != 0 {
		results["music"] = *musicRank
	}
	// Get place in the time streak leaderboard
	streakRank, err := globalsRepo.UserStatsRepository.GetUserLeaderboardRank(userId, "streak")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard streak rank for user %s", userId)
		return nil, err
	}
	if *streakRank != 0 {
		results["streak"] = *streakRank
	}

	return results, nil

}

func GiveTimeoutToMemberWithId(s *discordgo.Session, guildId string, userId string, reason string, creationTimestamp int64, sTimeoutLength float64) error {

	result := globalsRepo.TimeoutsRepository.GetTimeoutsCountForUser(userId)
	if result > 0 {
		return fmt.Errorf("a user cannot be given more than 1 timeout at a time")
	}

	// If the user is on their 10th timeout
	numArchivedTimeouts := globalsRepo.TimeoutsRepository.GetArchivedTimeoutsCountForUser(userId)
	if numArchivedTimeouts == 9 {
		// ban them instead
		err := s.GuildBanCreateWithReason(guildId, userId, "Received 10th and final timeout", 1)
		if err != nil {
			fmt.Println("Error banning user on 10th timeout: ", err)
			return err
		}
		// and clean DB related entries
		err = DeleteAllMemberData(userId)
		if err != nil {
			fmt.Println("Error deleting user data on 10th timeout: ", err)
			return err
		}
	}

	err := globalsRepo.TimeoutsRepository.SaveTimeout(userId, reason, creationTimestamp, int(sTimeoutLength))
	if err != nil {
		fmt.Printf("Error ocurred while storing timeout for user: %s\n", err)
		return fmt.Errorf(err.Error())
	}

	// Give actual Discord timeout to member
	timeoutExpiryTimestamp := time.Now().Add(time.Second * time.Duration(sTimeoutLength))
	err = s.GuildMemberTimeout(guildId, userId, &timeoutExpiryTimestamp)
	if err != nil {
		fmt.Println("Error timing out user: ", err)
		return fmt.Errorf("%v", err)
	}

	return nil

}

func SendDirectMessageToMember(s *discordgo.Session, userId string, msg string) error {
	errDm := dm.DmUser(s, userId, msg)
	if errDm != nil {
		fmt.Printf("Error sending DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}

func GetMemberTimeouts(userId string) (*dataModels.Timeout, []dataModels.ArchivedTimeout, error) {

	// Result variables
	var activeTimeoutResult *dataModels.Timeout = nil
	var archivedTimeoutResults []dataModels.ArchivedTimeout = []dataModels.ArchivedTimeout{}

	// Active timeout
	activeTimeout, err := globalsRepo.TimeoutsRepository.GetUserTimeout(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			activeTimeoutResult = nil
		} else {
			return nil, nil, err
		}
	}
	activeTimeoutResult = activeTimeout

	// Archived timeouts
	archivedTimeoutResults, err = globalsRepo.TimeoutsRepository.GetAllArchivedTimeoutsForUser(userId)
	if err != nil {
		return nil, nil, err
	}

	return activeTimeoutResult, archivedTimeoutResults, nil

}

func ClearMemberActiveTimeout(s *discordgo.Session, guildId string, userId string) error {

	err := globalsRepo.TimeoutsRepository.ClearTimeoutForUser(userId)
	if err != nil {
		return err
	}

	err = s.GuildMemberTimeout(guildId, userId, nil)
	if err != nil {
		fmt.Println("Error timing out user: ", err)
		return fmt.Errorf("%v", err)
	}

	return nil

}

func AppealTimeout(guildId string, userId string) error {

	activeTimeout, _, err := GetMemberTimeouts(userId)
	if err != nil {
		timeoutError := fmt.Errorf("an error ocurred while retrieving timeout data for user with ID %s: %v\n", userId, err)
		return timeoutError
	}

	if activeTimeout == nil {
		return fmt.Errorf("no active timeout was found for user with ID `%s`\n", userId)
	}

	// TODO

	return nil

}

func GetMemberExperiencePoints(userId string) (*float64, error) {

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving User from DB: %v\n", err)
		return nil, err
	}

	return &user.CurrentExperience, nil

}

func GrantMemberExperience(userId string, activityType string, points float64) (float64, error) {

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving user with UID %s from OTA DB: %v\n", userId, err)
		return -1, err
	}

	err = globalsRepo.UsersRepository.AddUserExpriencePoints(userId, points)
	if err != nil {
		fmt.Printf("An error ocurred while granting XP to user: %v\n", err)
		return -1, err
	}

	// Also store records for the monthly leaderboard
	monthlyEntryExists := globalsRepo.MonthlyLeaderboardRepository.EntryExists(userId)
	if monthlyEntryExists <= 0 {
		if monthlyEntryExists == -1 {
			return -1, fmt.Errorf("monthly leaderboard entry to was not found in the DB; likely an error has ocurred")
		}
		// Entry doesn't exist for member, so create one
		err := globalsRepo.MonthlyLeaderboardRepository.AddLeaderboardEntry(userId, user.Gender)
		if err != nil {
			return -1, err
		}
	}
	err = globalsRepo.MonthlyLeaderboardRepository.AddLeaderboardExpriencePoints(userId, points)
	if err != nil {
		fmt.Printf("An error ocurred while granting monthly leaderboard XP to user: %v\n", err)
		return -1, err
	}

	return user.CurrentExperience + points, nil

}

func RemoveMemberExperience(userId string, activityType string) (*float64, error) {

	isMember := globalsRepo.UsersRepository.UserExists(userId)
	if isMember <= 0 {
		if isMember == -1 {
			return nil, fmt.Errorf("member to grant XP to was not found in the DB; likely an error has ocurred")
		}
		return nil, fmt.Errorf("member to remove XP from was not found in the DB; likely the given member is a bot application")
	}

	var xpToRemove float64
	switch activityType {
	case "MSG_REWARD":
		xpToRemove = globals.ExperienceReward_MessageSent
	case "REACT_REWARD":
		xpToRemove = globals.ExperienceReward_ReactionReceived
	case "SLASH_REWARD":
		xpToRemove = globals.ExperienceReward_SlashCommandUsed
	case "IN_VC_REWARD":
		xpToRemove = globals.ExperienceReward_InVc
	case "IN_MUSIC_REWARD":
		xpToRemove = globals.ExperienceReward_InMusic
	}

	err := globalsRepo.UsersRepository.RemoveUserExpriencePoints(userId, xpToRemove)
	if err != nil {
		fmt.Printf("An error ocurred while removing XP from user: %v\n", err)
		return nil, err
	}

	// Also remove points from the monthly leaderboard
	monthlyEntryExists := globalsRepo.MonthlyLeaderboardRepository.EntryExists(userId)
	if monthlyEntryExists <= 0 {
		if monthlyEntryExists == -1 {
			return nil, fmt.Errorf("monthly leaderboard entry to was not found in the DB; likely an error has ocurred")
		}
	}

	if monthlyEntryExists == 1 {
		err = globalsRepo.MonthlyLeaderboardRepository.RemoveUserExpriencePoints(userId, xpToRemove)
		if err != nil {
			fmt.Printf("An error ocurred while removing monthly leaderboard XP from user: %v\n", err)
			return nil, err
		}
	}

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving User (%s) from DB after removing XP. Member may have left the server.\n", userId)
		return nil, err
	}

	return &user.CurrentExperience, nil

}

func IsBot(s *discordgo.Session, guildId string, userId string, debug bool) (*bool, error) {

	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		if debug {
			log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		}
		return nil, err
	}

	isBot := apiUser.Bot

	return &isBot, nil
}

func SetGender(userId string, genderValue string) error {

	user, err := globalsRepo.UsersRepository.GetUser(userId)
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

	_, err = globalsRepo.UsersRepository.UpdateUser(*user)
	if err != nil {
		return err
	}

	// Also set gender in leaderboard - if applicable
	count := globalsRepo.MonthlyLeaderboardRepository.EntryExists(userId)
	if count <= 0 {
		if count == -1 {
			return fmt.Errorf("an error ocurred while checking for user leaderboard entry")
		}
	} else {
		err = globalsRepo.MonthlyLeaderboardRepository.UpdateCategoryForUser(userId, user.Gender)
		if err != nil {
			return err
		}

	}

	return nil

}

func JailMember(s *discordgo.Session, guildId string, userId string, reason string, jailRoleName string, notificationChannelId string) (*dataModels.JailedUser, *dataModels.User, error) {

	var err error

	currentTimestamp := time.Now()

	// Pick a random task to assign to the jailed user
	taskToFree := utils.GetRandomFromArray(globals.JailTasks)

	// Build a record of the jailed user for the command feedback
	var jailedRecord *dataModels.JailedUser = &dataModels.JailedUser{
		UserId:         userId,
		JailedAt:       currentTimestamp.Unix(),
		TaskToComplete: taskToFree,
		Reason:         reason,
	}

	// Add User to Jail in the DB
	err = globalsRepo.JailRepository.AddUserToJail(userId, reason, taskToFree, currentTimestamp.Unix())
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Give designated Discord role to member to restrict access
	err = GiveDiscordRoleToMember(s, guildId, userId, jailRoleName)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Send notification about jailing on designated channel
	embed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️    A New Prisoner Has Arrived").
		AddField("Known As", user.DiscordTag, false).
		AddField("Convicted Because", reason, false).
		AddField("Tasked With", taskToFree, false).
		AddField("Convincted At", currentTimestamp.String(), false)

	go notifications.SendEmbedToTextChannel(s, notificationChannelId, *embed)

	// Send Jail DM to jailed user
	// TODO

	return jailedRecord, user, nil

}
