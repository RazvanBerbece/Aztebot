package repositories

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	databaseconn "github.com/LxrdVixxeN/azteca-discord/internal/bot-service/data/connection"
	dataModels "github.com/LxrdVixxeN/azteca-discord/internal/bot-service/data/models"
)

type UsersRepositoryInterface interface {
}

type UsersRepository struct {
	conn databaseconn.Database
}

func NewUsersRepository() *UsersRepository {
	repo := new(UsersRepository)
	repo.conn.ConnectDatabaseHandle()
	return repo
}

func (r UsersRepository) GetRolesForUser(userId string) ([]dataModels.Role, error) {

	// Get assigned role IDs for given user from the DB
	rows, err := r.conn.Db.Query("SELECT currentRoleIds FROM Users WHERE userId = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
	}
	defer rows.Close()

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var roleIdsString string
	var placeholders string
	var ids []int
	for rows.Next() {
		if err := rows.Scan(&roleIdsString); err != nil {
			return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
		}
		placeholders, ids = getSqlFriendlyListOfStringIDs(roleIdsString)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
	}
	// Check for zero rows - if the query arg only has the opening paranthesis
	if len(roleIdsString) == 0 {
		return nil, fmt.Errorf("GetRolesForUser %s - User: No roles found for user. User may not exist", userId)
	}

	// Get roles from DB with found role IDs and return a list of them
	roles, err := r.GetRolesByIds(placeholders, ids)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - Roles: %v", userId, err)
	}

	return roles, nil
}

func (r UsersRepository) GetRolesByIds(placeholders string, ids []int) ([]dataModels.Role, error) {

	// Convert roleIDIntegers to a slice of interface{} to use as variadic args in Db.Query()
	var rolesAsListOfAny []interface{}
	for _, id := range ids {
		rolesAsListOfAny = append(rolesAsListOfAny, id)
	}

	var roles []dataModels.Role
	query := fmt.Sprintf("SELECT * FROM Roles WHERE id IN (%s)", placeholders)
	rowsRoles, err := r.conn.Db.Query(query, rolesAsListOfAny...)
	if err != nil {
		return nil, fmt.Errorf("GetRolesByIds <%d>: %v", ids, err)
	}
	defer rowsRoles.Close()
	for rowsRoles.Next() {
		var role dataModels.Role
		if err := rowsRoles.Scan(&role.Id, &role.RoleName, &role.DisplayName, &role.Colour, &role.Info, &role.Perms); err != nil {
			return nil, fmt.Errorf("GetRolesByIds: %v", err)
		}
		roles = append(roles, role)
	}
	if err := rowsRoles.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesByIds: %v", err)
	}
	// Check for zero rows - if the query arg has no IDs retrieved from the Users table
	if len(roles) == 0 {
		return nil, fmt.Errorf("GetRolesByIds: No roles found for ids %d", ids)
	}
	return roles, nil
}

// Method that returns a list of placeholders (?) and a list of IDs to be used in a
// `Select * From T Where k in (...)` SQL query.
func getSqlFriendlyListOfStringIDs(roleIdsString string) (string, []int) {
	roles := strings.Split(roleIdsString, ",")
	var placeholders []string
	for range roles {
		placeholders = append(placeholders, "?")
	}
	var roleIdIntegers []int
	for _, id := range roles {
		roleID, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			log.Fatal(err)
		}
		roleIdIntegers = append(roleIdIntegers, roleID)
	}
	return strings.Join(placeholders, ","), roleIdIntegers
}
