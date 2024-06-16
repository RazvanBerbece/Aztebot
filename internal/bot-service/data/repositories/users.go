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

	// Get role IDs for given user
	rows, err := r.conn.Db.Query("SELECT currentRoleIds FROM Users WHERE userId = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to variable fields.
	var inRolesSqlList string
	var placeholders string
	var ids []any
	for rows.Next() {
		var roleIdsString string
		if err := rows.Scan(&roleIdsString); err != nil {
			return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
		}
		placeholders, ids = getSqlFriendlyListOfStringIDs(roleIdsString)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
	}
	// Check for zero rows - if the query arg only has the opening paranthesis
	if len(inRolesSqlList) == 1 {
		return nil, fmt.Errorf("GetRolesForUser %s - User: No roles found for user. User may not exist", userId)
	}

	// Get roles with found roleIds and return
	var roles []dataModels.Role

	query := fmt.Sprintf("SELECT * FROM Roles WHERE id IN (%s)", placeholders)
	rowsRoles, err := r.conn.Db.Query(query, ids...)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - Roles <%s>: %v", userId, inRolesSqlList, err)
	}
	defer rowsRoles.Close()
	for rowsRoles.Next() {
		var role dataModels.Role
		if err := rowsRoles.Scan(&role.Id, &role.RoleName, &role.DisplayName, &role.Colour, &role.Info, &role.Perms); err != nil {
			return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
		}
		roles = append(roles, role)
	}
	if err := rowsRoles.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s - User: %v", userId, err)
	}
	// Check for zero rows - if the query arg only has the opening paranthesis
	if len(inRolesSqlList) == 1 {
		return nil, fmt.Errorf("GetRolesForUser %s - User: No roles found for user. User may not exist", userId)
	}

	return roles, nil
}

func getSqlFriendlyListOfStringIDs(roleIdsString string) (string, []any) {
	// SQL inclusion queries are used to then get the Role details.
	// So we need to build the argument for the `in` SQL clause
	// (example: (1, 2, 3) as in SELECT * FROM table WHERE id in (1, 2, 3))
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
	// Convert roleIDIntegers to a slice of interface{} to use as variadic asrgs in Db.Query()
	var rolesAsListOfAny []interface{}
	for _, id := range roleIdIntegers {
		rolesAsListOfAny = append(rolesAsListOfAny, id)
	}

	return strings.Join(placeholders, ","), rolesAsListOfAny
}
