package member

import (
	"fmt"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

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
