package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type RolesRepositoryInterface interface {
}

type RolesRepository struct {
	conn databaseconn.Database
}

func NewRolesRepository() *RolesRepository {
	repo := new(RolesRepository)
	repo.conn.ConnectDatabaseHandle()
	return repo
}

func (r RolesRepository) GetRole(displayName string) (*dataModels.Role, error) {

	r.conn.ConnectDatabaseHandle()

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Roles WHERE displayName = ?"
	row := r.conn.Db.QueryRow(query, displayName)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var role dataModels.Role
	err := row.Scan(
		&role.Id,
		&role.RoleName,
		&role.DisplayName,
		&role.Emoji,
		&role.Info,
	)

	if err != nil {
		return nil, err
	}

	defer r.conn.Db.Close()

	return &role, nil
}
