package member

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
)

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
