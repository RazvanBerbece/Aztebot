package member

import (
	"fmt"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

func AwardFunds(userId string, funds float64) error {

	if funds < 0 || funds > 500000.0 {
		return fmt.Errorf("cannot award funds to user with ID `%s`, because the number of awarded `funds` (`%.2f`) is invalid", userId, funds)
	}

	err := globalRepositories.WalletsRepository.AddFundsToWalletForUser(userId, funds)
	if err != nil {
		fmt.Printf("An error ocurred while awarding funds to user %s: %v\n", userId, err)
		return err
	}

	return nil

}
