package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type RolesRepository struct {
	conn databaseconn.Database
}

func NewRolesRepository() *RolesRepository {
	repo := new(RolesRepository)
	repo.conn.ConnectDatabaseHandle()
	return repo
}

func (r RolesRepository) GetRole(roleDisplayName string) (*dataModels.Role, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Roles WHERE displayName = ?"
	row := r.conn.Db.QueryRow(query, roleDisplayName)

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

	return &role, nil
}
