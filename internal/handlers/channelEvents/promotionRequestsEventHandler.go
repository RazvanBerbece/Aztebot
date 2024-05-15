package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandlePromotionRequestEvents(s *discordgo.Session, orderRoleNames []string) {

	for xpEvent := range globalMessaging.PromotionRequestsChannel {

		var userGuildId = xpEvent.GuildId
		var userId = xpEvent.UserId
		var userXp = xpEvent.CurrentXp
		var userCurrentLevel = xpEvent.CurrentLevel
		var userNumberMessagesSent = xpEvent.MessagesSent
		var userTimeSpentInVc = xpEvent.TimeSpentInVc

		// Check current stats against progression table
		// Figure out the promoted role to be given
		var promotedLevel int = 0
		var promotedRoleName string = ""
		switch {
		// No order
		case userXp < 7500:
			// no promotion required
			continue
		// First order
		case userXp >= 7500 && userXp < 10000:
			if userNumberMessagesSent >= 1000 && userTimeSpentInVc >= 60*60*15 {
				promotedLevel = 1
				promotedRoleName = "🔗 Zelator"
			}
		case userXp >= 10000 && userXp < 15000:
			if userNumberMessagesSent >= 2500 && userTimeSpentInVc >= 60*60*20 {
				promotedLevel = 2
				promotedRoleName = "📖 Theoricus"
			}
		case userXp >= 15000 && userXp < 20000:
			if userNumberMessagesSent >= 3500 && userTimeSpentInVc >= 60*60*25 {
				promotedLevel = 3
				promotedRoleName = "🎩 Practicus"
			}
		case userXp >= 20000 && userXp < 30000:
			if userNumberMessagesSent >= 5000 && userTimeSpentInVc >= 60*60*30 {
				promotedLevel = 4
				promotedRoleName = "📿 Philosophus"
			}
		// Second order
		case userXp >= 30000 && userXp < 45000:
			if userNumberMessagesSent >= 12500 && userTimeSpentInVc >= 60*60*40 {
				promotedLevel = 5
				promotedRoleName = "🔮 Adeptus Minor"
			}
		case userXp >= 45000 && userXp < 50000:
			if userNumberMessagesSent >= 15000 && userTimeSpentInVc >= 60*60*45 {
				promotedLevel = 6
				promotedRoleName = "〽️ Adeptus Major"
			}
		case userXp >= 50000 && userXp < 100000:
			if userNumberMessagesSent >= 20000 && userTimeSpentInVc >= 60*60*50 {
				promotedLevel = 7
				promotedRoleName = "🧿 Adeptus Exemptus"
			}
		// Third order
		case userXp >= 100000 && userXp < 150000:
			if userNumberMessagesSent >= 35000 && userTimeSpentInVc >= 60*60*200 {
				promotedLevel = 8
				promotedRoleName = "☀️ Magister Templi"
			}
		case userXp >= 150000 && userXp < 200000:
			if userNumberMessagesSent >= 45000 && userTimeSpentInVc >= 60*60*250 {
				promotedLevel = 9
				promotedRoleName = "🧙🏼 Magus"
			}
		case userXp >= 200000:
			if userNumberMessagesSent >= 50000 && userTimeSpentInVc >= 60*60*300 {
				promotedLevel = 10
				promotedRoleName = "⚔️ Ipsissimus"
			}
		}

		// Promotion is available for current member
		if promotedLevel > userCurrentLevel {

			fmt.Println("Processing promotion...")

			// Give promoted level in DB
			err := globalRepositories.UsersRepository.SetLevel(userId, promotedLevel)
			if err != nil {
				fmt.Printf("Error ocurred while setting member level in DB: %v\n", err)
			}

			fmt.Printf("Level up ! (%d -> %d)\n", userCurrentLevel, promotedLevel)

			// Give promoted role in DB (and cleanup the old one)
			promotedRole, err := globalRepositories.RolesRepository.GetRole(promotedRoleName) // to append
			if err != nil {
				fmt.Printf("Error ocurred while reading role from DB: %v\n", err)
			}
			fmt.Println("Role to promote to: ", promotedRole)

			if promotedLevel == 1 {
				fmt.Println("First level up !")
				// no previous order role so no need to remove it, only append to list of IDs
				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error ocurred while appending role ID to member in DB: %v\n", err)
				}
			} else if promotedLevel > 1 {
				currentOrderRole, err := member.GetMemberOrderRole(userId, orderRoleNames) // to remove
				if err != nil {
					fmt.Printf("Error ocurred while reading member order role from DB: %v\n", err)
				}
				fmt.Println("Current order role to remove:", currentOrderRole.DisplayName)
				err = globalRepositories.UsersRepository.RemoveUserRoleWithId(userId, currentOrderRole.Id)
				if err != nil {
					fmt.Printf("Error ocurred while removing member role from DB: %v\n", err)
				}
				fmt.Println("Removed:", currentOrderRole.DisplayName)
				err = globalRepositories.UsersRepository.AppendUserRoleWithId(userId, promotedRole.Id)
				if err != nil {
					fmt.Printf("Error ocurred while appending role ID to member in DB: %v\n", err)
				}
			}

			// Get refreshed role IDs after processing
			user, err := globalRepositories.UsersRepository.GetUser(userId)
			if err != nil {
				fmt.Printf("Error ocurred while retrieving user and roles from DB: %v\n", err)
			}
			fmt.Println("UPDATED ROLES:", user.CurrentRoleIds)
			err = member.RefreshDiscordRolesWithIdForMember(s, userGuildId, userId, user.CurrentRoleIds)
			if err != nil {
				fmt.Printf("Error ocurred while refreshing member roles on-Discord: %v\n", err)
			}
		}

	}

}
