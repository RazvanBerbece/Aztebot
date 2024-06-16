package startup

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncLevelsAtStartup(s *discordgo.Session, guildId string, uids []string) {

	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	// For all users in the database
	fmt.Println("[STARTUP] Checkpoint Task SyncLevelsAtStartup() at", time.Now(), "-> Updating", len(uids), "levels")
	for _, uid := range uids {
		err := member.ProcessProgressionForMember(uid, guildId)
		if err != nil {
			fmt.Printf("Failed to process level at startup for member with uid %s: %v\n", uid, err)
			continue
		}
	}

	// Cleanup repos
	go utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task SyncLevelsAtStartup() at", time.Now())

}
