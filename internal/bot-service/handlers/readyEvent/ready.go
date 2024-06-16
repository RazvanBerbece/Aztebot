package readyEvent

import (
	"fmt"
	"time"

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
	s.UpdateGameStatus(0, "Type /help")

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
		for range ticker.C {
			// Run your periodic task here
			UpdateUsersInCron(s)
		}
	}()

}

func UpdateUsersInCron(s *discordgo.Session) error {

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globals.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("Error retrieving members:", err)
		return err
	}

	// Process the current batch of members
	processMembers(s, members)

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
		processMembers(s, members)
	}

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member) {
	// Your logic to process members goes here
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := utils.SyncUser(s, globals.DiscordMainGuildId, member.User.ID, member)
		if err != nil {
			fmt.Printf("Error syncinc member %s: %v", member.User.Username, err)
		}
	}
}
