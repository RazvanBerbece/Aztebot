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

	// Set initial status for the AzteBot
	s.UpdateGameStatus(0, "/help")

	// Other setups

	// Cron func to sync users and their DB entity
	var interval int
	if globals.UserSyncIntervalErr != nil {
		interval = 60
	} else {
		interval = globals.UserSyncInterval
	}
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		// Define the repositories here (and establish the connections) in order to not flood the DB with connection attempts
		rolesRepository := repositories.NewRolesRepository()
		usersRepository := repositories.NewUsersRepository()
		for range ticker.C {
			UpdateUsersInCron(s, rolesRepository, usersRepository)
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
