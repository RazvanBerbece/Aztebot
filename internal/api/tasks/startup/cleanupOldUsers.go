package startup

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func CleanupMemberAtStartup(s *discordgo.Session, uids []string) error {

	fmt.Println("[STARTUP] Starting Task CleanupMemberAtStartup() at", time.Now())

	// Inject new connections
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	uidsLength := len(uids)

	// For each tag in the DB, delete user from table
	for i := 0; i < uidsLength; i++ {
		uid := uids[i]
		_, err := s.GuildMember(globalConfiguration.DiscordMainGuildId, uid)
		if err != nil {
			// if the member does not exist on the main server, delete from the database
			err = member.DeleteAllMemberData(uid)
			if err != nil {
				fmt.Println("Error deleting hanging user data on startup sync: ", err)
			}
		}
	}

	// Cleanup repos
	go utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task CleanupMemberAtStartup() at", time.Now())

	return nil

}
