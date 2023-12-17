package memberUpdateEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func MemberRoleUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	logging.LogHandlerCall("MemberRoleUpdate", "")

	// Sync user in DB with the current Discord member state
	fmt.Printf("Handling role update for %s\n", m.Member.User.Username)
	err := utils.SyncUser(s, m.GuildID, m.Member.User.ID, m.Member)
	if err != nil {
		fmt.Printf("Error ocurred while syncing new user roles with the DB")
	}

}
