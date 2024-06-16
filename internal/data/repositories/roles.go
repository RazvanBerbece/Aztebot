package repositories

import (
	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
)

type RolesRepository struct {
	Conn databaseconn.Database
}

func NewRolesRepository() *RolesRepository {
	repo := new(RolesRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r RolesRepository) GetRole(roleDisplayName string) (*dataModels.Role, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Roles WHERE displayName = ?"
	row := r.Conn.Db.QueryRow(query, roleDisplayName)

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

func (r RolesRepository) GetRoleById(roleId int) (*dataModels.Role, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Roles WHERE id = ?"
	row := r.Conn.Db.QueryRow(query, roleId)

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
