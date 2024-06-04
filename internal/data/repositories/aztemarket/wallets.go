package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
)

type DbWalletsRepository interface {
	AwardFunds(userId string, funds float64)
}

type WalletsRepository struct {
	DbContext databaseconn.AztemarketDbContext
}

func NewWalletsRepository(connString string) WalletsRepository {
	repo := WalletsRepository{databaseconn.AztemarketDbContext{
		ConnectionString: connString,
	}}
	repo.DbContext.Connect()
	return repo
}

func (r WalletsRepository) AddFundsToWalletForUser(userId string, funds float64) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
		UPDATE Wallets SET 
			funds = funds + ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(funds, userId)
	if err != nil {
		return err
	}

	return nil
}
