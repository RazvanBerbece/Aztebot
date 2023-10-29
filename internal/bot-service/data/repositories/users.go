package repositories

import (
	"fmt"
	"strconv"
	"strings"

	dataModels "github.com/LxrdVixxeN/azteca-discord/internal/bot-service/data/models"
	databaseconn "github.com/RazvanBerbece/UnifyFootballBot/internal/database-conn"
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

	// Loop through rows, using Scan to assign column data to struct fields.
	var roleIds []int
	for rows.Next() {
		var roleIdsString string
		if err := rows.Scan(&roleIdsString); err != nil {
			return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
		}

		// Split string on whitespace to get a lsit of int role IDs
		roles := strings.Split(roleIdsString, " ")
		for _, roleId := range roles {
			idAsInt, err := strconv.Atoi(roleId)
			if err != nil {
				return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
			}
			roleIds = append(roleIds, int(idAsInt))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
	}
	// Check for zero rows
	if len(roleIds) == 0 {
		return nil, fmt.Errorf("GetRolesForUser %s: No roles found for user. User may not exist.", userId)
	}

	// Get roles with found roleIds and return
	var roles []dataModels.Role

	rowsRoles, err := r.conn.Db.Query("SELECT * FROM Roles WHERE roleId in = ?", roleIds)
	if err != nil {
		return nil, fmt.Errorf("GetRolesForUser %s: %v", userId, err)
	}
	defer rows.Close()

	return roles, nil
}
