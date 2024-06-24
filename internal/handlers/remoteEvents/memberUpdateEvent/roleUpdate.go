package memberUpdateEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/logging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func MemberRoleUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	// Only run this handler if it's a role update event or if we can't check for the previous member state
	if m.BeforeUpdate != nil {
		if utils.EqualSlices(m.BeforeUpdate.Roles, m.Roles) {
			// no change in roles, return early
			return
		} else {
			if globalConfiguration.AuditRoleUpdatesInChannel { // Log the extended role update (previous and current roles, etc.)
				go generateExtendedRoleUpdateLog(s, m)
			}
		}
	}

	// Sync user in DB with the current Discord member state
	go member.SyncMember(s, m.GuildID, m.Member.User.ID, m.Member, globalConfiguration.OrderRoleNames, false)
}

func generateExtendedRoleUpdateLog(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// Get previous roles
	var previousRoles = m.BeforeUpdate.Roles
	var previousRolesString = ""
	for idx, roleId := range previousRoles {
		role, err := member.GetDiscordRole(s, m.GuildID, roleId)
		if err != nil {
			fmt.Printf("Error ocurred while retrieving Discord role: %v\n", err)
		}
		if idx < len(previousRoles)-1 {
			previousRolesString += fmt.Sprintf("`%s`,", role.Name)
		} else if idx == len(previousRoles)-1 {
			previousRolesString += fmt.Sprintf("`%s`", role.Name)
		}
	}
	if len(previousRolesString) == 0 {
		previousRolesString = "_none found in the state cache._"
	}

	// Get current member roles
	currentRoles := m.Roles
	currentRolesString := ""
	for idx, roleId := range currentRoles {
		role, err := member.GetDiscordRole(s, m.GuildID, roleId)
		if err != nil {
			fmt.Printf("Error ocurred while retrieving Discord role: %v\n", err)
		}
		if idx < len(currentRoles)-1 {
			currentRolesString += fmt.Sprintf("`%s`,", role.Name)
		} else if idx == len(currentRoles)-1 {
			currentRolesString += fmt.Sprintf("`%s`", role.Name)
		}
	}

	if !utils.EqualSlices(previousRoles, currentRoles) {
		// Audit update by logging in provided debug channel
		logMsg := fmt.Sprintf("Handling role update for `%s` [`%s`]\n\nPrevious roles: %s\n\nUpdated roles: %s", m.Member.User.Username, m.Member.User.ID, previousRolesString, currentRolesString)
		discordChannelLogger := logging.NewDiscordLogger(s, "notif-debug")
		go discordChannelLogger.LogInfo(logMsg)
	}

}
