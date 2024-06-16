package repositories

import (
	"fmt"
	"strconv"
	"strings"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
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

func (r UsersRepository) GetUser(userId string) (*dataModels.User, error) {
	// Get assigned role IDs for given user from the DB
	query := "SELECT * FROM Users WHERE userId = ?"
	row := r.conn.Db.QueryRow(query, userId)

	// Scan the role IDs and process them into query arguments to use
	// in the Roles table
	var user dataModels.User
	err := row.Scan(&user.Id,
		&user.DiscordTag,
		&user.UserId,
		&user.CurrentRoleIds,
		&user.CurrentCircle,
		&user.CurrentInnerOrder,
		&user.CurrentLevel,
		&user.CurrentExperience,
		&user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("GetUser %s: %v", userId, err)
	}

	return &user, nil
}

func (r UsersRepository) SaveInitialUserDetails(tag string, userId string) (*dataModels.User, error) {

	user := &dataModels.User{
		DiscordTag:        tag,
		UserId:            userId,
		CurrentRoleIds:    "",
		CurrentCircle:     "",
		CurrentInnerOrder: nil,
		CurrentLevel:      0,
		CurrentExperience: 0,
		CreatedAt:         nil,
	}

	stmt, err := r.conn.Db.Prepare(`
		INSERT INTO 
			Users(discordTag, 
				userId, 
				currentRoleIds, 
				currentCircle, 
				currentInnerOrder, 
				currentLevel, 
				currentExperience, 
				createdAt) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.DiscordTag, user.UserId, user.CurrentRoleIds, user.CurrentCircle, user.CurrentInnerOrder, user.CurrentLevel, user.CurrentExperience, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r UsersRepository) UpdateUser() {}

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
		placeholders, ids = GetSqlPlaceholderAndValueRoleCommand(idArray(roleIdsString))
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
		if err := rowsRoles.Scan(&role.Id, &role.RoleName, &role.DisplayName, &role.Emoji, &role.Info); err != nil {
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
func GetSqlPlaceholderAndValueRoleCommand(roles []int) (string, []int) {
	var placeholders []string
	for range roles {
		placeholders = append(placeholders, "?")
	}
	return strings.Join(placeholders, ","), roles
}

func idArray(idsString string) []int {

	var ids []int
	stringIds := strings.Split(idsString, ",")

	for _, id := range stringIds {
		num, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Could not parse role ID %s into integer: %v", id, err)
			continue
		}
		ids = append(ids, num)
	}

	return ids

}
