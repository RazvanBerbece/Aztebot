package startup

import (
	"database/sql"
	"fmt"
	"time"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncExperiencePointsGainsAtStartup(s *discordgo.Session, uids []string) {

	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

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

		currentXp := user.CurrentExperience

		stats, errStats := userStatsRepository.GetStatsForUser(uid)
		if errStats != nil {
			if errStats == sql.ErrNoRows {
				continue
			}
			fmt.Println("[STARTUP] Failed Task SyncExperiencePointsGainsAtStartup() at", time.Now(), "for UID", "with error", errStats)
		}

		computedXp := utils.CalculateExperiencePointsFromStats(
			stats.NumberMessagesSent,
			stats.NumberSlashCommandsUsed,
			stats.NumberReactionsReceived,
			stats.TimeSpentInVoiceChannels,
			stats.TimeSpentListeningToMusic,
			globalConfiguration.DefaultExperienceReward_MessageSent,
			globalConfiguration.DefaultExperienceReward_SlashCommandUsed,
			globalConfiguration.DefaultExperienceReward_ReactionReceived,
			globalConfiguration.DefaultExperienceReward_InVc,
			globalConfiguration.DefaultExperienceReward_InMusic)

		// XP can also be added via bonuses, so take that into consideration
		// by picking the max between them
		// todo: this could instead use an XpGrants table ?
		var actualXp float64
		if computedXp != currentXp {
			// mismatch resolution
			if currentXp > computedXp {
				actualXp = currentXp // current XP would include awards, grants, etc.
			} else {
				actualXp = computedXp
			}
		} else {
			actualXp = computedXp
		}
		user.CurrentExperience = actualXp

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
