package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type RolesRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewRolesRepository() *RolesRepository {
	repo := new(RolesRepository)
	repo.Conn.Connect()
	return repo
}

func (r RolesRepository) GetRole(roleDisplayName string) (*dax.Role, error) {

	query := "SELECT * FROM Roles WHERE displayName = ?"
	row := r.Conn.SqlDb.QueryRow(query, roleDisplayName)

	var role dax.Role
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

func (r RolesRepository) GetRoleById(roleId int) (*dax.Role, error) {

	query := "SELECT * FROM Roles WHERE id = ?"
	row := r.Conn.SqlDb.QueryRow(query, roleId)

	var role dax.Role
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

func (r RolesRepository) RoleByIdExists(roleId int) (*dax.Role, error) {

	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Roles WHERE id = ?"
	row := r.Conn.SqlDb.QueryRow(query, roleId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var role dax.Role
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

func (r RolesRepository) RoleByDisplayNameExists(roleDisplayName string) int {
	query := "SELECT COUNT(*) FROM Roles WHERE displayName = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, roleDisplayName).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for role in OTA DB: %v\n", err)
		return -1
	}
	return count
}
