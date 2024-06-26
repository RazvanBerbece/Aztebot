package startup

import (
	"fmt"
	"time"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	memberService "github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SyncMembersAtStartup(s *discordgo.Session, defaultOrderRoleNames []string, syncProgression bool) error {

	fmt.Println("[STARTUP] Starting Task SyncMembersAtStartup() at", time.Now())

	// Inject new connections
	rolesRepository := repositories.NewRolesRepository()
	usersRepository := repositories.NewUsersRepository()
	userStatsRepository := repositories.NewUsersStatsRepository()
	jailRepository := repositories.NewJailRepository()

	// Retrieve all members in the guild
	members, err := s.GuildMembers(globalConfiguration.DiscordMainGuildId, "", 1000)
	if err != nil {
		fmt.Println("[STARTUP] Failed Task SyncMembersAtStartup() at", time.Now(), "with error", err)
		return err
	}

	// Process the current batch of members
	processMembers(s, members, rolesRepository, usersRepository, userStatsRepository, jailRepository, defaultOrderRoleNames, syncProgression)

	// Paginate
	for len(members) != 0 {
		// Set the 'After' parameter to the ID of the last member in the current batch
		lastMemberID := members[len(members)-1].User.ID

		members, err = s.GuildMembers(globalConfiguration.DiscordMainGuildId, lastMemberID, 1000)
		if err != nil {
			fmt.Println("[STARTUP] Failed Task SyncMembersAtStartup() at", time.Now(), "with error", err)
			return err
		}

		// Process the next batch of members
		processMembers(s, members, rolesRepository, usersRepository, userStatsRepository, jailRepository, defaultOrderRoleNames, syncProgression)
	}

	// Cleanup
	go utils.CleanupRepositories(rolesRepository, usersRepository, userStatsRepository, nil, nil)

	fmt.Println("[STARTUP] Finished Task SyncMembersAtStartup() at", time.Now())

	return nil

}

func processMembers(s *discordgo.Session, members []*discordgo.Member, rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository, jailRepository *repositories.JailRepository, defaultOrderRoleNames []string, syncProgression bool) {
	for _, member := range members {
		// If it's a bot, skip
		if member.User.Bot {
			continue
		}
		// For each member, sync their details (either add to DB or update)
		err := memberService.SyncMemberPersistent(s, globalConfiguration.DiscordMainGuildId, member.User.ID, member, rolesRepository, usersRepository, userStatsRepository, jailRepository, defaultOrderRoleNames, syncProgression)
		if err != nil {
			fmt.Printf("Error ocurred while processing batch of members for ID %s: %v", member.User.ID, err)
		}
	}
}
