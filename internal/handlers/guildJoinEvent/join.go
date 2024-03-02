package joinEvent

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a new joined member.
func GuildJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	// If it's a bot, skip
	if m.Member.User.Bot {
		return
	}

	utils.LogHandlerCall("GuildJoin", "")

	// Store newly-joined user to DB tables (probably only the initial details and awaiting for verification and cron sync)
	err := utils.SyncUser(s, globals.DiscordMainGuildId, m.Member.User.ID, m.Member)
	if err != nil {
		fmt.Printf("Error storing new member %s to DB: %v", m.Member.User.Username, err)
	}

	// Create entity for user stats in the DB
	err = globalsRepo.UserStatsRepository.SaveInitialUserStats(m.Member.User.ID)
	if err != nil {
		fmt.Printf("Error storing new member stats %s to DB: %v", m.Member.User.Username, err)
	}

	// Other actions to do on guild join
	// e.g guide DM, etc.

}
