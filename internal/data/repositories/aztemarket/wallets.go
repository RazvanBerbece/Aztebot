package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
)

type DbWalletsRepository interface {
	AddFundsToWalletForUser(userId string, funds float64) error
	GetWalletIdForUser(userId string) (*string, error)
	DeleteWalletForUser(userId string) error
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

func (r WalletsRepository) GetWalletIdForUser(userId string) (*string, error) {

	query := "SELECT id FROM Wallets WHERE userId = ?"
	row := r.DbContext.SqlDb.QueryRow(query, userId)

	var id string
	err := row.Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil

}

func (r WalletsRepository) DeleteWalletForUser(userId string) error {

	query := "DELETE FROM Wallets WHERE userId = ?"

	_, err := r.DbContext.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting wallet entry for user: %w", err)
	}

	return nil

}
