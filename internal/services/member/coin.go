package member

import (
	"fmt"

	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

func AwardFunds(userId string, funds int) error {

	if funds < 0 || funds > 500000 {
		return fmt.Errorf("cannot award funds to user with ID `%s`, because the number of awarded `funds` (`%d`) is invalid", userId, funds)
	}

	err := globalRepositories.WalletsRepository.AddFundsToWalletForUser(userId, funds)
	if err != nil {
		return err
	}

	return nil

}
