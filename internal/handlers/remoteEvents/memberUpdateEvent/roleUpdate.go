package memberUpdateEvent

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func MemberRoleUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	// DEBUG
	// TODO: Remove this once role-related issues have been identified
	currentRoles := m.Roles
	currentRolesString := ""
	for _, roleId := range currentRoles {
		role, err := member.GetDiscordRole(s, m.GuildID, roleId)
		if err != nil {
			fmt.Printf("Error ocurred while retrieving Discord role: %v\n", err)
		}
		currentRolesString += fmt.Sprintf("%s,", role.Name)
	}
	fmt.Printf("Handling role update for %s (updated roles: %s)\n", m.Member.User.Username, currentRolesString)

	// Sync user in DB with the current Discord member state
	err := member.SyncMember(s, m.GuildID, m.Member.User.ID, m.Member, globalConfiguration.OrderRoleNames, globalConfiguration.SyncProgressionInMemberUpdates)
	if err != nil {
		fmt.Printf("Error ocurred while syncing new user roles with the DB: %v\n", err)
	}
}
