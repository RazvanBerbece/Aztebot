package repositories

import (
	"fmt"
	"strconv"
	"strings"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type UsersRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewUsersRepository() *UsersRepository {
	repo := new(UsersRepository)
	repo.Conn.Connect()
	return repo
}

func (r UsersRepository) GetAllDiscordUids() ([]string, error) {

	var userIds []string
	rowsUsers, err := r.Conn.SqlDb.Query("SELECT userId FROM Users")
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

	return userIds, nil
}

func (r UsersRepository) GetAllUsers() ([]dax.User, error) {

	var users []dax.User

	rowsUsers, err := r.Conn.SqlDb.Query("SELECT * FROM Users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}

	defer rowsUsers.Close()

	for rowsUsers.Next() {
		var user dax.User
		if err := rowsUsers.Scan(&user.Id,
			&user.DiscordTag,
			&user.UserId,
			&user.CurrentRoleIds,
			&user.CurrentCircle,
			&user.CurrentInnerOrder,
			&user.CurrentLevel,
			&user.CurrentExperience,
			&user.CreatedAt,
			&user.Gender); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, user)
	}
	if err := rowsUsers.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}

	return users, nil
}

func (r UsersRepository) UserExists(userId string) int {
	query := "SELECT COUNT(*) FROM Users WHERE userId = ?"
	var count int
	err := r.Conn.SqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		fmt.Printf("An error ocurred while checking for user in OTA DB: %v\n", err)
		return -1
	}
	return count
}

func (r UsersRepository) GetUser(userId string) (*dax.User, error) {

	query := "SELECT * FROM Users WHERE userId = ?"
	row := r.Conn.SqlDb.QueryRow(query, userId)

	var user dax.User
	err := row.Scan(&user.Id,
		&user.DiscordTag,
		&user.UserId,
		&user.CurrentRoleIds,
		&user.CurrentCircle,
		&user.CurrentInnerOrder,
		&user.CurrentLevel,
		&user.CurrentExperience,
		&user.CreatedAt,
		&user.Gender)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UsersRepository) DeleteUser(userId string) error {

	query := "DELETE FROM Users WHERE userId = ?"

	_, err := r.Conn.SqlDb.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (r UsersRepository) SaveInitialUserDetails(tag string, userId string, timestamp *int64) (*dax.User, error) {

	user := &dax.User{
		DiscordTag:        tag,
		UserId:            userId,
		CurrentRoleIds:    "",
		CurrentCircle:     "",
		CurrentInnerOrder: nil,
		CurrentLevel:      0,
		CurrentExperience: 0,
		CreatedAt:         timestamp,
		Gender:            -1,
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
		INSERT INTO 
			Users(
				discordTag, 
				userId, 
				currentRoleIds, 
				currentCircle, 
				currentInnerOrder, 
				currentLevel, 
				currentExperience, 
				createdAt,
				gender
			)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.DiscordTag, user.UserId, user.CurrentRoleIds, user.CurrentCircle, user.CurrentInnerOrder, user.CurrentLevel, user.CurrentExperience, user.CreatedAt, user.Gender)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r UsersRepository) UpdateUser(user dax.User) (*dax.User, error) {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			discordTag = ?, 
			currentRoleIds = ?, 
			currentCircle = ?, 
			currentInnerOrder = ?, 
			currentLevel = ?, 
			currentExperience = ?, 
			createdAt = ?,
			gender = ?
		WHERE userId = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.DiscordTag, user.CurrentRoleIds, user.CurrentCircle, user.CurrentInnerOrder, user.CurrentLevel, user.CurrentExperience, user.CreatedAt, user.Gender, user.UserId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UsersRepository) SetUserCreatedAt(userId string, timestamp int64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			createdAt = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(timestamp, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) AddUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
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

func (r UsersRepository) SetLevel(userId string, level int) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentLevel = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(level, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) RemoveUserExpriencePoints(userId string, experiencePoints float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
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

func (r UsersRepository) RemoveUserRoles(userId string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentRoleIds = ","
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) AppendUserRoleWithId(userId string, roleId int) error {

	roles, err := r.GetRolesForUser(userId)
	if err != nil {
		return err
	}

	roleIdsString := ""
	for _, role := range roles {
		roleIdsString += fmt.Sprintf("%d,", role.Id)
	}
	roleIdsString += fmt.Sprintf("%d,", roleId)

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentRoleIds = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleIdsString, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) RemoveUserRoleWithId(userId string, roleId int) error {

	roles, err := r.GetRolesForUser(userId)
	if err != nil {
		return err
	}

	roleIdsString := ""
	for _, role := range roles {
		if role.Id != roleId {
			roleIdsString += fmt.Sprintf("%d,", role.Id)
		}
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentRoleIds = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleIdsString, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) RemoveUserRoleWithName(userId string, roleDisplayName string) error {

	roles, err := r.GetRolesForUser(userId)
	if err != nil {
		return err
	}

	roleIdsString := ""
	for _, role := range roles {
		if role.DisplayName != roleDisplayName {
			roleIdsString += fmt.Sprintf("%d,", role.Id)
		}
	}

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentRoleIds = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleIdsString, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) SetUserRoles(userId string, roleIdsString string) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE Users SET 
			currentRoleIds = ?
		WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, roleIdsString)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) GetRolesForUser(userId string) ([]dax.Role, error) {

	// Get assigned role IDs for given user from the DB
	rows, err := r.Conn.SqlDb.Query("SELECT currentRoleIds FROM Users WHERE userId = ?", userId)
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
		if roleIdsString == "" {
			// no roles assigned for user
			return []dax.Role{}, nil
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

func (r UsersRepository) GetRolesByIds(placeholders string, ids []int) ([]dax.Role, error) {

	// Convert roleIDIntegers to a slice of interface{} to use as variadic args in Db.Query()
	var rolesAsListOfAny []interface{}
	for _, id := range ids {
		rolesAsListOfAny = append(rolesAsListOfAny, id)
	}

	// Retrieve the roles in ascending order of importance (higher id means higher importance)
	var roles []dax.Role
	query := fmt.Sprintf("SELECT * FROM Roles WHERE id IN (%s) ORDER BY id ASC", placeholders)
	rowsRoles, err := r.Conn.SqlDb.Query(query, rolesAsListOfAny...)
	if err != nil {
		return nil, fmt.Errorf("GetRolesByIds <%d>: %v", ids, err)
	}
	defer rowsRoles.Close()
	for rowsRoles.Next() {
		var role dax.Role
		if err := rowsRoles.Scan(&role.Id, &role.RoleName, &role.DisplayName, &role.Emoji, &role.Info); err != nil {
			return nil, fmt.Errorf("GetRolesByIds: %v", err)
		}
		roles = append(roles, role)
	}
	if err := rowsRoles.Err(); err != nil {
		return nil, fmt.Errorf("GetRolesByIds: %v", err)
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
