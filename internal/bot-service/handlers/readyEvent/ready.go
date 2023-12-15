package readyEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

var ()

// Called once the Discord servers confirm a succesful connection.
func Ready(s *discordgo.Session, event *discordgo.Ready) {

	logging.LogHandlerCall("Ready", "")

	// Define the repositories here for the cron functions (and reuse their connections)
	// in order to not flood the DB with connection attempts
	rolesRepository := repositories.NewRolesRepository()
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Other setups

	// Cron funcs to sync users and their DB entity
	var syncInterval int
	var cleanupInterval int
	if globals.UserSyncIntervalErr != nil {
		fmt.Printf("Could not parse UserSyncInterval environment variable: %v\n", globals.UserSyncIntervalErr)
		syncInterval = 60
	} else {
		syncInterval = globals.UserSyncInterval
	}
	if globals.UserCleanupIntervalErr != nil {
		fmt.Printf("Could not parse UserCleanupInterval environment variable: %v\n", globals.UserCleanupIntervalErr)
		cleanupInterval = 60
	} else {
		cleanupInterval = globals.UserCleanupInterval
	}

	syncTicker := time.NewTicker(time.Second * time.Duration(syncInterval))
	cleanupTicker := time.NewTicker(time.Second * time.Duration(cleanupInterval))

	// Periodic sync of the members on the server with the DB
	go func() {
		for range syncTicker.C {
			UpdateUsersInCron(s, rolesRepository, usersRepository, userStatsRepository)
		}
	}()

	// Periodic cleanup of users from the DB
	go func() {
		for range cleanupTicker.C {
			CleanupUsersInCron(s, usersRepository, userStatsRepository)
		}
	}()

}

func UpdateUsersInCron(s *discordgo.Session, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) error {

	fmt.Println("Starting Task UpdateUsersInCron() at", time.Now())

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globals.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("Failed Task UpdateUsersInCron() at", time.Now(), "with error", err)
		return err
	}

	// Process the current batch of members
	processMembers(s, members, rolesRepository, usersRepository, userStatsRepository)

	// Paginate
	for len(members) == 1000 {
		// Set the 'After' parameter to the ID of the last member in the current batch
		lastMemberID := members[len(members)-1].User.ID
		members, err = s.GuildMembers(globals.DiscordMainGuildId, lastMemberID, 1000)
		if err != nil {
			fmt.Println("Failed Task UpdateUsersInCron() at", time.Now(), "with error", err)
			return err
		}

		// Process the next batch of members
		processMembers(s, members, rolesRepository, usersRepository, userStatsRepository)
	}

	fmt.Println("Finished Task UpdateUsersInCron() at", time.Now())

	return nil

}

func CleanupUsersInCron(s *discordgo.Session, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) error {

	fmt.Println("Starting Task CleanupUsersInCron() at", time.Now())

	// Retrieve all members from the DB
	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("Failed Task CleanupUsersInCron() at", time.Now(), "with error", err)
		return err
	}

	// For each tag in the DB, delete user from table
	for _, uid := range uids {
		_, err := s.GuildMember(globals.DiscordMainGuildId, uid)
		if err != nil {
			// if the member does not exist on the main server, delete from the database
			// delete user stats
			err := userStatsRepository.DeleteUserStats(uid)
			if err != nil {
				fmt.Println("Failed Task CleanupUsersInCron() at", time.Now(), "with error", err)
				return err
			}
			// delete user
			errUsers := usersRepository.DeleteUser(uid)
			if errUsers != nil {
				fmt.Println("Failed Task CleanupUsersInCron() at", time.Now(), "with error", errUsers)
				return errUsers
			}
		}
		// Sleep for a bit to not exceed request frequency limits for the Discord API
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Println("Finished Task CleanupUsersInCron() at", time.Now())

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) {
	// Your logic to process members goes here
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := utils.SyncUserPersistent(s, globals.DiscordMainGuildId, member.User.ID, member, rolesRepository, usersRepository, userStatsRepository)
		if err != nil {
			fmt.Printf("Error syncinc member %s: %v", member.User.Username, err)
		}
	}
}
