package member

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

func GetMemberExperiencePoints(userId string) (*float64, error) {

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving User from DB: %v\n", err)
		return nil, err
	}

	return &user.CurrentExperience, nil

}

func GetMemberXpRank(userId string) (*int, error) {

	xpRank, err := globalRepositories.UserStatsRepository.GetUserXpRank(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard XP rank for user %s", userId)
		return nil, err
	}

	return xpRank, nil
}

func GetMemberRankInLeaderboards(userId string) (map[string]int, error) {

	results := make(map[string]int)

	// Get place in the messages sent leaderboard
	msgRank, err := globalRepositories.UserStatsRepository.GetUserLeaderboardRank(userId, "msg")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard msg rank for user %s", userId)
		return nil, err
	}
	if *msgRank != 0 {
		results["msg"] = *msgRank
	}
	// Get place in the reactions received leaderboard
	reactRank, err := globalRepositories.UserStatsRepository.GetUserLeaderboardRank(userId, "react")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard react rank for user %s", userId)
		return nil, err
	}
	if *reactRank != 0 {
		results["react"] = *reactRank
	}
	// Get place in the time spent in VCs leaderboard
	vcRank, err := globalRepositories.UserStatsRepository.GetUserLeaderboardRank(userId, "vc")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard vc rank for user %s", userId)
		return nil, err
	}
	if *vcRank != 0 {
		results["vc"] = *vcRank
	}
	// Get place in the time spent in music channels leaderboard
	musicRank, err := globalRepositories.UserStatsRepository.GetUserLeaderboardRank(userId, "music")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard music rank for user %s", userId)
		return nil, err
	}
	if *musicRank != 0 {
		results["music"] = *musicRank
	}
	// Get place in the time streak leaderboard
	streakRank, err := globalRepositories.UserStatsRepository.GetUserLeaderboardRank(userId, "streak")
	if err != nil {
		fmt.Printf("An error ocurred while retrieving leaderboard streak rank for user %s", userId)
		return nil, err
	}
	if *streakRank != 0 {
		results["streak"] = *streakRank
	}

	return results, nil

}

func GrantMemberExperience(userId string, points float64) (float64, error) {

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving user with UID %s from OTA DB: %v\n", userId, err)
		return -1, err
	}

	err = globalRepositories.UsersRepository.AddUserExpriencePoints(userId, points)
	if err != nil {
		fmt.Printf("An error ocurred while granting XP to user: %v\n", err)
		return -1, err
	}

	// If not a staff membber
	if !IsStaff(userId, globalConfiguration.StaffRoles) {
		// Also store records for the monthly leaderboard
		monthlyEntryExists := globalRepositories.MonthlyLeaderboardRepository.EntryExists(userId)
		if monthlyEntryExists <= 0 {
			if monthlyEntryExists == -1 {
				return -1, fmt.Errorf("monthly leaderboard entry to was not found in the DB; likely an error has ocurred")
			}
			// Entry doesn't exist for member, so create one
			err := globalRepositories.MonthlyLeaderboardRepository.AddLeaderboardEntry(userId, user.Gender)
			if err != nil {
				return -1, err
			}
		}
		err = globalRepositories.MonthlyLeaderboardRepository.AddLeaderboardExpriencePoints(userId, points)
		if err != nil {
			fmt.Printf("An error ocurred while granting monthly leaderboard XP to user: %v\n", err)
			return -1, err
		}

		// Also store records for the daily leaderboard
		dailyEntryExists := globalRepositories.DailyLeaderboardRepository.EntryExists(userId)
		if dailyEntryExists <= 0 {
			if dailyEntryExists == -1 {
				return -1, fmt.Errorf("daily leaderboard entry to was not found in the DB; likely an error has ocurred")
			}
			// Entry doesn't exist for member, so create one
			err := globalRepositories.DailyLeaderboardRepository.AddLeaderboardEntry(userId, user.Gender)
			if err != nil {
				return -1, err
			}
		}
		err = globalRepositories.DailyLeaderboardRepository.AddLeaderboardExpriencePoints(userId, points)
		if err != nil {
			fmt.Printf("An error ocurred while granting daily leaderboard XP to user: %v\n", err)
			return -1, err
		}
	}

	// Kickstart automatic progression process
	err = ProcessProgressionForMember(userId, globalConfiguration.DiscordMainGuildId)
	if err != nil {
		fmt.Printf("An error ocurred while starting automatic progression for user: %v\n", err)
		return -1, err
	}

	return user.CurrentExperience + points, nil

}

func RemoveMemberExperience(userId string, activityType string) (*float64, error) {

	isMember := globalRepositories.UsersRepository.UserExists(userId)
	if isMember <= 0 {
		if isMember == -1 {
			return nil, fmt.Errorf("member to grant XP to was not found in the DB; likely an error has ocurred")
		}
		return nil, fmt.Errorf("member to remove XP from was not found in the DB; likely the given member is a bot application")
	}

	var xpToRemove float64
	switch activityType {
	case "MSG_REWARD":
		xpToRemove = globalConfiguration.ExperienceReward_MessageSent
	case "REACT_REWARD":
		xpToRemove = globalConfiguration.ExperienceReward_ReactionReceived
	case "SLASH_REWARD":
		xpToRemove = globalConfiguration.ExperienceReward_SlashCommandUsed
	case "IN_VC_REWARD":
		xpToRemove = globalConfiguration.ExperienceReward_InVc
	case "IN_MUSIC_REWARD":
		xpToRemove = globalConfiguration.ExperienceReward_InMusic
	}

	err := globalRepositories.UsersRepository.RemoveUserExpriencePoints(userId, xpToRemove)
	if err != nil {
		fmt.Printf("An error ocurred while removing XP from user: %v\n", err)
		return nil, err
	}

	// Also remove points from the monthly and daily leaderboards
	if !IsStaff(userId, globalConfiguration.StaffRoles) {
		monthlyEntryExists := globalRepositories.MonthlyLeaderboardRepository.EntryExists(userId)
		if monthlyEntryExists <= 0 {
			if monthlyEntryExists == -1 {
				return nil, fmt.Errorf("monthly leaderboard entry to was not found in the DB; likely an error has ocurred")
			}
		}

		if monthlyEntryExists == 1 {
			err = globalRepositories.MonthlyLeaderboardRepository.RemoveUserExpriencePoints(userId, xpToRemove)
			if err != nil {
				fmt.Printf("An error ocurred while removing monthly leaderboard XP from user: %v\n", err)
				return nil, err
			}
		}

		dailyEntryExists := globalRepositories.DailyLeaderboardRepository.EntryExists(userId)
		if dailyEntryExists <= 0 {
			if dailyEntryExists == -1 {
				return nil, fmt.Errorf("daily leaderboard entry to was not found in the DB; likely an error has ocurred")
			}
		}

		if dailyEntryExists == 1 {
			err = globalRepositories.DailyLeaderboardRepository.RemoveUserExpriencePoints(userId, xpToRemove)
			if err != nil {
				fmt.Printf("An error ocurred while removing daily leaderboard XP from user: %v\n", err)
				return nil, err
			}
		}
	}

	// Kickstart automatic progression process
	err = ProcessProgressionForMember(userId, globalConfiguration.DiscordMainGuildId)
	if err != nil {
		fmt.Printf("An error ocurred while starting automatic progression for user: %v\n", err)
		return nil, err
	}

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("An error ocurred while retrieving User (%s) from DB after removing XP. Member may have left the server.\n", userId)
		return nil, err
	}

	return &user.CurrentExperience, nil

}
