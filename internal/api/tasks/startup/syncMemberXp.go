package startup

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncExperiencePointsGainsAtStartup(s *discordgo.Session) {

	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "with error", err)
	}

	// For all users in the database
	fmt.Println("[STARTUP] Checkpoint Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "-> Updating", len(uids), "XP gains")
	for _, uid := range uids {

		user, err := usersRepository.GetUser(uid)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", err)
		}

		stats, errStats := userStatsRepository.GetStatsForUser(uid)
		if errStats != nil {
			if errStats == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", errStats)
		}

		updatedXp := utils.CalculateExperiencePointsFromStats(
			stats.NumberMessagesSent,
			stats.NumberSlashCommandsUsed,
			stats.NumberReactionsReceived,
			stats.TimeSpentInVoiceChannels,
			stats.TimeSpentListeningToMusic,
			globals.ExperienceReward_MessageSent,
			globals.ExperienceReward_SlashCommandUsed,
			globals.ExperienceReward_ReactionReceived,
			globals.ExperienceReward_InVc,
			globals.ExperienceReward_InMusic)
		user.CurrentExperience = float64(updatedXp)

		// Update user entity with new XP value
		_, err = usersRepository.UpdateUser(*user)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", err)
		}

	}

	// Cleanup repos
	go utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task SyncExperiencePointsGainsAtStartup() at", time.Now())

}
