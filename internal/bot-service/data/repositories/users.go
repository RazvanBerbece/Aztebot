package repositories

import (
	"fmt"
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
	rows, err := r.conn.Db.Query("SELECT [currentRoleIds] FROM Users WHERE userId = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to variable fields.
	var inRolesSqlList = "("
	for rows.Next() {
		var roleIdsString string
		if err := rows.Scan(&roleIdsString); err != nil {
			return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
		}

		// SQL inclusion queries are used to then get the Role details. So we need to build the argument for the `in` SQL clause.
		// example: SELECT * FROM table WHERE id in (1, 2, 3)
		roles := strings.Split(roleIdsString, " ")
		for index, roleId := range roles {
			if index == len(roles)-1 {
				inRolesSqlList = inRolesSqlList + fmt.Sprintf("\"%s\")", roleId)
				break
			}
			inRolesSqlList = inRolesSqlList + fmt.Sprintf("\"%s\",", roleId)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
	}
	// Check for zero rows - if the query arg only has the opening paranthesis
	if len(inRolesSqlList) == 1 {
		return nil, fmt.Errorf("GetRolesForUser %s: No roles found for user. User may not exist.", userId)
	}

	// Get roles with found roleIds and return
	var roles []dataModels.Role

	rowsRoles, err := r.conn.Db.Query("SELECT * FROM Roles WHERE roleId in ?", inRolesSqlList)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
	}
	defer rowsRoles.Close()

	return roles, nil
}
