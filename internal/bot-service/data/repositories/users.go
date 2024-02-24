package repositories

import (
	"fmt"
	"strconv"
	"strings"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/connection"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
)

type UsersRepository struct {
	Conn databaseconn.Database
}

func NewUsersRepository() *UsersRepository {
	repo := new(UsersRepository)
	repo.Conn.ConnectDatabaseHandle()
	return repo
}

func (r UsersRepository) GetAllDiscordUids() ([]string, error) {

	var userIds []string
	rowsUsers, err := r.Conn.Db.Query("SELECT userId FROM Users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUids: %v", err)
	}
	defer rowsUsers.Close()
	for rowsUsers.Next() {
		var id string
		if err := rowsUsers.Scan(&id); err != nil {
			return nil, fmt.Errorf("GetAllUids: %v", err)
		}
		userIds = append(userIds, id)
	}
	if err := rowsUsers.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUids: %v", err)
	}
	// Check for zero rows - if the query arg has no IDs retrieved from the Users table
	if len(userIds) == 0 {
		return nil, fmt.Errorf("GetAllUids: No users found in Users table")
	}

	if err != nil {
		return nil, err
	}

	return userIds, nil
}

func (r UsersRepository) GetAllUsers() ([]dataModels.User, error) {

	var users []dataModels.User

	rowsUsers, err := r.Conn.Db.Query("SELECT * FROM Users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}

	defer rowsUsers.Close()

	for rowsUsers.Next() {
		var user dataModels.User
		if err := rowsUsers.Scan(&user.Id,
			&user.DiscordTag,
			&user.UserId,
			&user.CurrentRoleIds,
			&user.CurrentCircle,
			&user.CurrentInnerOrder,
			&user.CurrentLevel,
			&user.CurrentExperience,
			&user.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, user)
	}
	if err := rowsUsers.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	// Check for zero rows - if the query arg has no IDs retrieved from the Users table
	if len(users) == 0 {
		return nil, fmt.Errorf("GetAllUsers: No users found in Users table")
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r UsersRepository) UserExists(userId string) int {
	query := "SELECT COUNT(*) FROM Users WHERE userId = ?"
	var count int
	err := r.Conn.Db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func (r UsersRepository) GetUser(userId string) (*dataModels.User, error) {

	query := "SELECT * FROM Users WHERE userId = ?"
	row := r.Conn.Db.QueryRow(query, userId)

	var user dataModels.User
	err := row.Scan(&user.Id,
		&user.DiscordTag,
		&user.UserId,
		&user.CurrentRoleIds,
		&user.CurrentCircle,
		&user.CurrentInnerOrder,
		&user.CurrentLevel,
		&user.CurrentExperience,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UsersRepository) DeleteUser(userId string) error {

	query := "DELETE FROM Users WHERE userId = ?"

	_, err := r.Conn.Db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
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

	stmt, err := r.Conn.Db.Prepare(`
		INSERT INTO 
			Users(
				discordTag, 
				userId, 
				currentRoleIds, 
				currentCircle, 
				currentInnerOrder, 
				currentLevel, 
				currentExperience, 
				createdAt
			)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?);`)
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

func (r UsersRepository) UpdateUser(user dataModels.User) (*dataModels.User, error) {

	stmt, err := r.Conn.Db.Prepare(`
		UPDATE Users SET 
			discordTag = ?, 
			currentRoleIds = ?, 
			currentCircle = ?, 
			currentInnerOrder = ?, 
			currentLevel = ?, 
			currentExperience = ?, 
			createdAt = ? 
		WHERE userId = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.DiscordTag, user.CurrentRoleIds, user.CurrentCircle, user.CurrentInnerOrder, user.CurrentLevel, user.CurrentExperience, user.CreatedAt, user.UserId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UsersRepository) AddUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.Db.Prepare(`
		UPDATE Users SET 
			currentExperience = currentExperience + ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(experiencePoints, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) RemoveUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.Db.Prepare(`
		UPDATE Users SET 
			currentExperience = currentExperience - ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(experiencePoints, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) GetRolesForUser(userId string) ([]dataModels.Role, error) {

	// Get assigned role IDs for given user from the DB
	rows, err := r.Conn.Db.Query("SELECT currentRoleIds FROM Users WHERE userId = ?", userId)
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

	// Retrieve the roles in ascending order of importance (higher id means higher importance)
	var roles []dataModels.Role
	query := fmt.Sprintf("SELECT * FROM Roles WHERE id IN (%s) ORDER BY id ASC", placeholders)
	rowsRoles, err := r.Conn.Db.Query(query, rolesAsListOfAny...)
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
		if id == "" {
			continue
		}
		num, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Could not parse role ID %s into integer: %v", id, err)
			continue
		}
		ids = append(ids, num)
	}

	return ids

}
