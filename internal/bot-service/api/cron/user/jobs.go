package cronUser

import (
	"fmt"
	"sync"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncUsersAtStartup(s *discordgo.Session) error {

	fmt.Println("Starting Task SyncUsersAtStartup() at", time.Now())

	// Inject new connections
	rolesRepository := repositories.NewRolesRepository()
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globals.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("Failed Task SyncUsersAtStartup() at", time.Now(), "with error", err)
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
			fmt.Println("Failed Task SyncUsersAtStartup() at", time.Now(), "with error", err)
			return err
		}

		// Process the next batch of members
		processMembers(s, members, rolesRepository, usersRepository, userStatsRepository)
	}

	// Cleanup
	utils.CleanupRepositories(rolesRepository, usersRepository, userStatsRepository, nil)

	fmt.Println("Finished Task SyncUsersAtStartup() at", time.Now())

	return nil

}

func CleanupMemberAtStartup(s *discordgo.Session, uids []string) error {

	fmt.Println("Starting Task CleanupMemberAtStartup() at", time.Now())

	// Inject new connections
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()

	uidsLength := len(uids)

	// For each tag in the DB, delete user from table
	var wg sync.WaitGroup
	wg.Add(uidsLength)
	for i := 0; i < uidsLength; i++ {
		go func(i int) {
			defer wg.Done()
			uid := uids[i]
			_, err := s.GuildMember(globals.DiscordMainGuildId, uid)
			if err != nil {
				// if the member does not exist on the main server, delete from the database
				// delete user stats
				err := userStatsRepository.DeleteUserStats(uid)
				if err != nil {
					fmt.Println("Failed Task CleanupMemberAtStartup() at", time.Now(), "with error", err)
					return
				}
				// delete user
				errUsers := usersRepository.DeleteUser(uid)
				if errUsers != nil {
					fmt.Println("Failed Task CleanupMemberAtStartup() at", time.Now(), "with error", errUsers)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	// Cleanup
	utils.CleanupRepositories(nil, usersRepository, userStatsRepository, nil)

	fmt.Println("Finished Task CleanupMemberAtStartup() at", time.Now())

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository) {
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := utils.SyncUserPersistent(s, globals.DiscordMainGuildId, member.User.ID, member, rolesRepository, usersRepository, userStatsRepository)
		if err != nil && err.Error() != "no update was executed" {
			fmt.Printf("Error syncing member %s: %v\n", member.User.Username, err)
		}
	}
}
