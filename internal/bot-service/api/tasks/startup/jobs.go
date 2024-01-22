package startup

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
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

func SendInformationEmbedsToTextChannels(s *discordgo.Session) {

	var textChannels map[string]string

	// TODO: Make the channels and their descriptions use environment variables somehow
	if globals.Environment == "staging" {
		// Dev text channels
		textChannels = map[string]string{
			"1188135110042734613": "default",
			"1194451477192773773": "staff-rules",
			"1198686819928264784": "server-rules",
		}
	} else {
		// Production text channels
		textChannels = map[string]string{
			"1176277764001767464": "info-music",
			"1100486860058398770": "staff-rules",
			"1100142572141281460": "server-rules",
		}
	}

	for id, details := range textChannels {
		hasMessage, err := utils.ChannelHasDefaultInformationMessage(s, id)
		if err != nil {
			fmt.Printf("Could not check for default message in channel %s (%s): %v", id, details, err)
			continue
		}
		if hasMessage {
			// Do not send this default message as it already exists
			continue
		} else {
			// Send associated default message to given text channel
			var embedText string
			var hasOwnEmbed bool
			var ownEmbed *embed.Embed = embed.NewEmbed()
			switch details {
			case "default":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/default.txt")
			case "info-music":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/music-info.txt")
			case "staff-rules":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/staff-rules.txt")
				hasOwnEmbed = true
				mutateLongEmbedFromStaticData(embedText, ownEmbed)
			case "server-rules":
				embedText = utils.GetTextFromFile("internal/bot-service/handlers/readyEvent/assets/defaultContent/server-rules.txt")
				hasOwnEmbed = true
				mutateLongEmbedFromStaticData(embedText, ownEmbed)
			}

			var messageEmbedToPost *discordgo.MessageEmbed
			if !hasOwnEmbed {
				messageEmbedToPost = embed.NewEmbed().
					SetTitle("🤖  Information Message").
					SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
					SetColor(000000).
					AddField("", embedText, false).
					MessageEmbed
			} else {
				messageEmbedToPost = ownEmbed.MessageEmbed
			}

			_, err := s.ChannelMessageSendEmbed(id, messageEmbedToPost)
			if err != nil {
				log.Fatalf("An error occured while sending a default message (%s): %v", details, err)
				return
			}
		}
	}

}

// Note that this is a mutating function on `hasOwnEmbed` and `embed`.
func mutateLongEmbedFromStaticData(embedText string, embed *embed.Embed) {
	// Split the content into sections based on double newline characters ("\n\n")
	sections := strings.Split(embedText, "\n\n")
	for _, section := range sections {
		lines := strings.Split(section, "\n")
		if len(lines) > 0 {
			// Use the first line as the title and the rest as content
			title := lines[0]
			content := strings.Join(lines[1:], "\n")
			embed.AddField(title, content, false)
		}
	}
}
