package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztemarket"
)

type DbCurrencySystemStateRepository interface {
	CreateCurrencySystem(guildId string,
		currencyName string,
		totalCurrencyAvailable float64,
		totalCurrencyInFlow float64,
		dateOfLastReplenish int64,
	) (*dax.CurrencySystemState, error)
	GetCurrencyStateForGuild(guildId string) (*dax.CurrencySystemState, error)
	ReplenishCurrencyForGuild(guildId string, currencyAmount float64) error
	DeallocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error
	AllocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error
	DeleteCurrencySystem(guildId string) error
}

type CurrencySystemStateRepository struct {
	DbContext databaseconn.AztemarketDbContext
}

func NewCurrencySystemStateRepository(connString string) CurrencySystemStateRepository {
	repo := CurrencySystemStateRepository{databaseconn.AztemarketDbContext{
		ConnectionString: connString,
	}}
	repo.DbContext.Connect()
	return repo
}

func (r CurrencySystemStateRepository) CreateCurrencySystem(guildId string,
	currencyName string,
	totalCurrencyAvailable float64,
	totalCurrencyInFlow float64,
	dateOfLastReplenish int64,
) (*dax.CurrencySystemState, error) {

	css := &dax.CurrencySystemState{
		GuildId:                guildId,
		CurrencyName:           currencyName,
		TotalCurrencyAvailable: totalCurrencyAvailable,
		TotalCurrencyInFlow:    totalCurrencyInFlow,
		DateOfLastReplenish:    dateOfLastReplenish,
	}

	stmt, err := r.DbContext.SqlDb.Prepare(`
		INSERT INTO 
			CurrencySystemState(
				guildId, 
				currencyName, 
				totalCurrencyAvailable,
				totalCurrencyInFlow,
				dateOfLastReplenish
			)
		VALUES(?, ?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(css.GuildId, css.CurrencyName, css.TotalCurrencyAvailable, css.TotalCurrencyInFlow, css.DateOfLastReplenish)
	if err != nil {
		return nil, err
	}

	return css, nil

}

func (r CurrencySystemStateRepository) DeleteCurrencySystem(guildId string) error {

	query := "DELETE FROM CurrencySystemState WHERE guildId = ?"

	_, err := r.DbContext.SqlDb.Exec(query, guildId)
	if err != nil {
		return fmt.Errorf("error deleting wallet entry for user: %w", err)
	}

	return nil

}

func (r CurrencySystemStateRepository) ReplenishCurrencyForGuild(guildId string, currencyAmount float64) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
	UPDATE CurrencySystemState SET 
		totalCurrencyAvailable = totalCurrencyAvailable + ?
	WHERE guildId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(currencyAmount, guildId)
	if err != nil {
		return err
	}

	return nil

}

func (r CurrencySystemStateRepository) GetCurrencyStateForGuild(guildId string) (*dax.CurrencySystemState, error) {

	query := "SELECT * FROM CurrencySystemState WHERE guildId = ?"
	row := r.DbContext.SqlDb.QueryRow(query, guildId)

	var css dax.CurrencySystemState
	err := row.Scan(&css.GuildId,
		&css.CurrencyName,
		&css.TotalCurrencyAvailable,
		&css.TotalCurrencyInFlow,
		&css.DateOfLastReplenish)

	if err != nil {
		return nil, err
	}

	return &css, nil

}

func (r CurrencySystemStateRepository) AllocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
	UPDATE CurrencySystemState SET 
		totalCurrencyAvailable = totalCurrencyAvailable - ?,
		totalCurrencyInFlow = totalCurrencyInFlow + ?
	WHERE guildId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(currencyAmount, currencyAmount, guildId)
	if err != nil {
		return err
	}

	return nil
}

func (r CurrencySystemStateRepository) DeallocateFlowingCurrencyForGuild(guildId string, currencyAmount float64) error {
	stmt, err := r.DbContext.SqlDb.Prepare(`
	UPDATE CurrencySystemState SET 
		totalCurrencyInFlow = totalCurrencyInFlow - ?,
		totalCurrencyAvailable = totalCurrencyAvailable + ?
	WHERE guildId = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(currencyAmount, currencyAmount, guildId)
	if err != nil {
		return err
	}

	return nil
}
