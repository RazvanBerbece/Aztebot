package member

import (
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func ReplyComplexToInteraction(s *discordgo.Session, i *discordgo.Interaction, embed embed.Embed, actionsRow discordgo.ActionsRow, pageSize int) error {

	originalAllFields := make([]*discordgo.MessageEmbedField, len(embed.Fields))
	copy(originalAllFields, embed.Fields)

	pages := (len(originalAllFields) + pageSize - 1) / pageSize
	if len(embed.Fields) < pageSize {
		_, err := s.ChannelMessageSendEmbed(i.ChannelID, embed.MessageEmbed)
		if err != nil {
			return err
		}
	} else {
		embed.Fields = embed.Fields[0:pageSize]
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page 1 / %d", pages),
		}

		editContent := ""
		editWebhook := discordgo.WebhookEdit{
			Content:    &editContent,
			Embeds:     &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			Components: &[]discordgo.MessageComponent{actionsRow},
		}
		msg, err := s.InteractionResponseEdit(i, &editWebhook)
		if err != nil {
			fmt.Printf("Error sending complex response to interaction: %v\n", err)
			return err
		}

		// Keep paginated embeds in-memory to enable handling on button presses
		globalState.EmbedsToPaginate[msg.ID] = dataModels.EmbedData{
			ChannelId:   msg.ChannelID,
			FieldData:   &originalAllFields, // all fields
			CurrentPage: 1,                  // same for all complex paginated embeds
			Timestamp:   float64(time.Now().Unix()),
		}
	}

	return nil
}

func GetDiscordRole(s *discordgo.Session, guildId string, roleId string) (*discordgo.Role, error) {

	roles, err := s.GuildRoles(guildId)
	if err != nil {
		fmt.Printf("Error retrieving roles for guild with ID %s: %v\n", guildId, err)
		return nil, err
	}

	for _, role := range roles {
		if role.ID == roleId {
			return role, nil
		}
	}

	return nil, fmt.Errorf("a role with ID %s hasn't been found", roleId)

}

func IsBot(s *discordgo.Session, guildId string, userId string, debug bool) (*bool, error) {

	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		if debug {
			log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		}
		return nil, err
	}

	isBot := apiUser.Bot

	return &isBot, nil
}

// Removes all roles on the actual Discord member.
func RemoveAllDiscordRolesFromMember(s *discordgo.Session, guildId string, userId string) error {

	// Get the member's roles
	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return err
	}

	// Find all user's roles and delete them
	for _, roleID := range member.Roles {

		// 20 Mar 2024: Discord does not allow any way of removing the default Server Booster role from a guild member
		// so we just ignore it like it doesn't exist and hope that it goes away. :thumbs_down
		role, err := GetDiscordRole(s, guildId, roleID)
		if err != nil {
			fmt.Printf("Error retrieving role with ID %s: %v\n", roleID, err)
			return err
		}
		if role.Name == globalConfiguration.ServerBoosterDefaultRoleName {
			continue
		}

		err = s.GuildMemberRoleRemove(guildId, userId, roleID)
		if err != nil {
			fmt.Printf("Error removing role with ID %s: %v\n", roleID, err)
			return err
		}
	}

	return nil

}

func RemoveDiscordRoleFromMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// 20 Mar 2024: Same Server Booster trick as above
	if roleName == globalConfiguration.ServerBoosterDefaultRoleName {
		return nil
	}

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Remove the role by role ID from the Discord member
	err := s.GuildMemberRoleRemove(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Println("Error removing role from Discord member:\n", err)
		return err
	}

	return nil

}

func AddDiscordRoleToMember(s *discordgo.Session, guildId string, userId string, roleName string) error {

	// 20 Mar 2024: Same Server Booster trick as above
	if roleName == globalConfiguration.ServerBoosterDefaultRoleName {
		return nil
	}

	// Get the ID of the given role by name
	discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, roleName)
	if discordRoleId == nil {
		return fmt.Errorf("%s Discord Role not found in Guild", roleName)
	}

	// Add the role by role ID to the Discord member
	err := s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
	if err != nil {
		fmt.Printf("Error giving role with name %s to Discord member: %v\n", roleName, err)
		return err
	}

	// If a staff role, add a default 'STAFF Team' role as well
	if utils.StringInSlice(roleName, globalConfiguration.StaffRoles) {
		discordDefaultStaffRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, "STAFF Team")
		if discordDefaultStaffRoleId != nil {
			err = s.GuildMemberRoleAdd(guildId, userId, *discordDefaultStaffRoleId)
			if err != nil {
				fmt.Println("Error adding default STAFF role to Discord member:", err)
				return err
			}
		}
	}

	return nil

}

func AddDiscordRolesToMember(s *discordgo.Session, guildId string, userId string, roleIds []int) error {

	// For each role
	for _, roleId := range roleIds {
		role, err := globalRepositories.RolesRepository.GetRoleById(roleId)
		if err != nil {
			fmt.Printf("Error ocurred while adding DB roles to Discord member: %v\n", err)
			return err
		}

		// 20 Mar 2024: Same Server Booster trick as above
		if role.DisplayName == globalConfiguration.ServerBoosterDefaultRoleName {
			continue
		}

		// Get the role ID by display name from Discord
		discordRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, role.DisplayName)
		if discordRoleId != nil {
			// Add the role by role ID to the Discord member
			err = s.GuildMemberRoleAdd(guildId, userId, *discordRoleId)
			if err != nil {
				fmt.Printf("Error adding DB role with name %s to Discord member: %v\n", role.DisplayName, err)
				return err
			}

			// If a staff role, add a default 'STAFF Team' role as well
			if utils.StringInSlice(role.DisplayName, globalConfiguration.StaffRoles) {
				discordDefaultStaffRoleId := GetDiscordRoleIdForRoleWithName(s, guildId, "STAFF Team")
				if discordDefaultStaffRoleId != nil {
					err = s.GuildMemberRoleAdd(guildId, userId, *discordDefaultStaffRoleId)
					if err != nil {
						fmt.Println("Error adding default STAFF role to Discord member:", err)
						return err
					}
				}
			}
		}
	}

	// Process ORDER role post-update (based on the current role state)
	_, currentOrder := utils.GetCircleAndOrderForGivenRoles(roleIds)
	if currentOrder != nil {
		var discordOrderRoleIdToAdd *string
		if *currentOrder == 3 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
		} else if *currentOrder == 2 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
		} else if *currentOrder == 1 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
		}
		err := s.GuildMemberRoleAdd(guildId, userId, *discordOrderRoleIdToAdd)
		if err != nil {
			fmt.Println("Error adding order role to member:", err)
			return err
		}
	}

	return nil

}

func GetDiscordRolesForMember(s *discordgo.Session, guildId string, userId string) ([]discordgo.Role, error) {

	roles := []discordgo.Role{}

	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return nil, err
	}

	for _, roleID := range member.Roles {
		role, err := GetDiscordRole(s, guildId, roleID)
		if err != nil {
			fmt.Printf("Error retrieving role with ID %s: %v\n", roleID, err)
			return nil, err
		}
		roles = append(roles, *role)
	}

	return roles, nil
}

func GetDiscordOrderRoleNameForMember(s *discordgo.Session, guildId string, userId string) (*string, error) {

	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return nil, err
	}

	for _, roleID := range member.Roles {
		role, err := GetDiscordRole(s, guildId, roleID)
		if err != nil {
			fmt.Printf("Error retrieving role with ID %s: %v\n", roleID, err)
			return nil, err
		}
		if role.Name == "---- Third Order ----" || role.Name == "---- Second Order ----" || role.Name == "---- First Order ----" {
			return &role.Name, nil
		}
	}

	return nil, nil
}

// Recalculates and re-assigns the order Discord role for a member based on their role IDs in the DB.
func RefreshDiscordOrderRoleForMember(s *discordgo.Session, guildId string, userId string) error {

	// Retrieve current order role ID from the Discord servers
	// i.e ---- Third Order ----
	roles, err := GetDiscordRolesForMember(s, guildId, userId)
	if err != nil {
		fmt.Printf("Error retrieving Discord roles for member with UID %s: %v\n", userId, err)
		return err
	}
	for _, role := range roles {
		if role.Name == "---- Third Order ----" || role.Name == "---- Second Order ----" || role.Name == "---- First Order ----" {
			// And remove it on-Discord if ID is available
			err := s.GuildMemberRoleRemove(guildId, userId, role.ID)
			if err != nil {
				fmt.Println("Error removing order role from Discord member:\n", err)
				return err
			}
			break
		}
	}

	// Process ORDER role from the DB entry and assign in to the target member
	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		return err
	}
	var currentRoleIds = user.CurrentRoleIds
	_, currentOrder := utils.GetCircleAndOrderForGivenRoles(utils.GetRoleIdsFromRoleString(currentRoleIds))
	if currentOrder != nil {
		var discordOrderRoleIdToAdd *string
		if *currentOrder == 3 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Third Order ----")
		} else if *currentOrder == 2 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- Second Order ----")
		} else if *currentOrder == 1 {
			discordOrderRoleIdToAdd = GetDiscordRoleIdForRoleWithName(s, guildId, "---- First Order ----")
		}
		err := s.GuildMemberRoleAdd(guildId, userId, *discordOrderRoleIdToAdd)
		if err != nil {
			fmt.Println("Error adding order role to member:", err)
			return err
		}
	}

	return nil

}

func GetDiscordRoleIdForRoleWithName(s *discordgo.Session, guildId string, roleName string) *string {

	// Get the Guild
	guild, err := s.Guild(guildId)
	if err != nil {
		fmt.Println("Error retrieving guild:", err)
		return nil
	}

	// Find the Role ID based on the role's display name
	var roleID string = ""
	for _, role := range guild.Roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		fmt.Println("No role ID was found for role name", roleName)
		return nil
	}

	return &roleID
}
