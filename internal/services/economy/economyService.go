package economy

import (
	"fmt"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztemarket"
	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztemarket"
)

type EconomyService struct {
	// repos
	CurrencySystemStateRepositoryRepository repositories.DbCurrencySystemStateRepositoryRepository
}

func (s EconomyService) CreateCurrencySystem(guildId string, currencyName string, totalCurrencyAvailable float64, totalCurrencyInFlow float64, dateOfLastReplenish int64) (*dax.CurrencySystemState, error) {

	currencySystem, err := s.CurrencySystemStateRepositoryRepository.CreateCurrencySystem(guildId, currencyName, totalCurrencyAvailable, totalCurrencyInFlow, dateOfLastReplenish)
	if err != nil {
		return nil, fmt.Errorf("failed to create currency system for guild `%s`", guildId)
	}

	return currencySystem, nil
}

func (s EconomyService) GetCurrencyStateForGuild(guildId string) (*dax.CurrencySystemState, error) {

	currencySystem, err := s.CurrencySystemStateRepositoryRepository.GetCurrencyStateForGuild(guildId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve currency system state for guild `%s`", guildId)
	}

	return currencySystem, nil
}

func (s EconomyService) ReplenishCurrencyForGuild(guildId string, currencyAmount float64) error {

	err := s.CurrencySystemStateRepositoryRepository.ReplenishCurrencyForGuild(guildId, currencyAmount)
	if err != nil {
		return fmt.Errorf("failed to replenish currency for guild `%s`", guildId)
	}

	return nil
}

func (s EconomyService) AllocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error {

	err := s.CurrencySystemStateRepositoryRepository.AllocateFlowingCurrencyForGuild(guildId, currencyAmount)
	if err != nil {
		return fmt.Errorf("failed to allocate flowing currency for guild `%s`: %v", guildId, err)
	}

	return nil
}

func (s EconomyService) DeallocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error {

	err := s.CurrencySystemStateRepositoryRepository.DeallocateFlowingCurrencyForGuild(guildId, currencyAmount)
	if err != nil {
		return fmt.Errorf("failed to deallocate flowing currency for guild `%s`", guildId)
	}

	return nil
}

func (s EconomyService) DeleteCurrencySystem(guildId string) error {

	err := s.CurrencySystemStateRepositoryRepository.DeleteCurrencySystem(guildId)
	if err != nil {
		return fmt.Errorf("failed to delete currency system for guild `%s`", guildId)
	}

	return nil
}
