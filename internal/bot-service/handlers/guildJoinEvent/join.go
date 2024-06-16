package joinEvent

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/logging"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirms a new joined member.
func GuildJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	logging.LogHandlerCall("GuildJoin", "")

	// Register user details and initial roles into DB
	// usersRepository := repositories.NewUsersRepository()
	// user, err := usersRepository.GetUser(userId)
	// if err != nil {
	// 	log.Fatalf("Cannot retrieve user with id %s: %v", userId, err)
	// }

}
