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

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Other setups

	// Cron funcs to sync users and their DB entity
	var syncInterval int
	var cleanupInterval int
	if globals.UserSyncIntervalErr != nil {
		syncInterval = 60
	} else {
		syncInterval = globals.UserSyncInterval
	}
	if globals.UserCleanupIntervalErr != nil {
		cleanupInterval = 60
	} else {
		cleanupInterval = globals.UserCleanupInterval
	}

	syncTicker := time.NewTicker(time.Second * time.Duration(syncInterval))
	cleanupTicker := time.NewTicker(time.Second * time.Duration(cleanupInterval))

	// Periodic sync of the members on the server with the DB
	go func() {
		for range syncTicker.C {
			UpdateUsersInCron(s, rolesRepository, usersRepository)
		}
	}()

	// Periodic cleanup of users from the DB
	go func() {
		for range cleanupTicker.C {
			CleanupUsersInCron(s, usersRepository)
		}
	}()

}

func UpdateUsersInCron(s *discordgo.Session, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository) error {

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globals.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("Error retrieving members:", err)
		return err
	}

	// Process the current batch of members
	processMembers(s, members, rolesRepository, usersRepository)

	// Paginate
	for len(members) == 1000 {
		// Set the 'After' parameter to the ID of the last member in the current batch
		lastMemberID := members[len(members)-1].User.ID
		members, err = s.GuildMembers(globals.DiscordMainGuildId, lastMemberID, 1000)
		if err != nil {
			fmt.Println("Error retrieving members:", err)
			return err
		}

		// Process the next batch of members
		processMembers(s, members, rolesRepository, usersRepository)
	}

	fmt.Println("Ran Task UpdateUsersInCron() at", time.Now())

	return nil

}

func CleanupUsersInCron(s *discordgo.Session, usersRepository *repositories.UsersRepository) error {

	// Retrieve all members from the DB
	uids, err := usersRepository.GetAllDiscordUids()
	if err != nil {
		fmt.Println("Error retrieving user IDs from DB:", err)
		return err
	}

	// For each tag in the DB, delete user from table
	for _, uid := range uids {
		_, err := s.GuildMember(globals.DiscordMainGuildId, uid)
		if err != nil {
			// if the member does not exist on the main server
			if discordgoErr, ok := err.(*discordgo.RESTError); ok && discordgoErr.Message.Code == discordgo.ErrCodeUnknownMember {
				// delete from the database
				err := usersRepository.DeleteUser(uid)
				if err != nil {
					fmt.Printf("Error deleting left user with UID %s from DB: %v", uid, err)
					return err
				}
			} else {
				fmt.Println("Error retrieving member:", err)
				return err
			}
		}
	}

	fmt.Println("Ran Task CleanupUsersInCron() at", time.Now())

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository) {
	// Your logic to process members goes here
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := utils.SyncUserPersistent(s, globals.DiscordMainGuildId, member.User.ID, member, rolesRepository, usersRepository)
		if err != nil {
			fmt.Printf("Error syncinc member %s: %v", member.User.Username, err)
		}
	}
}
