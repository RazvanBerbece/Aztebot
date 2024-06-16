package slashHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashSync(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := ProcessUserUpdate(i.Interaction.Member.User.ID, s, i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to sync your data.",
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully synced data with the internal records.",
		},
	})
}

func ProcessUserUpdate(userId string, s *discordgo.Session, event *discordgo.InteractionCreate) error {

	rolesRepository := repositories.NewRolesRepository()
	usersRepository := repositories.NewUsersRepository()
	user, err := usersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
		if err == sql.ErrNoRows {
			log.Printf("Storing user with id %s", userId)
			user, err = usersRepository.SaveInitialUserDetails(event.Member.User.Username, userId)
			if err != nil {
				log.Fatalf("Cannot store user %s with id %s: %v", event.Member.User.Username, userId, err)
				return err
			}
		}
	}

	if user != nil {
		// Roles
		// Get current roles from user (as they appear on the Discord guild)
		var currentRoleIds string
		var roleIds []int
		for _, role := range event.Member.Roles {
			// Build a list of roles taken from the Discord guild
			// and then use the list to update the role IDs, circle and order in the database
			userRoleObj, err := s.State.Role(event.GuildID, role) // role DisplayName in OTA DB
			if err != nil {
				log.Println("Error getting role from Discord servers:", err)
				return err
			}
			roleDax, err := rolesRepository.GetRole(userRoleObj.Name)
			if err != nil {
				log.Println("Error getting role from DB:", err)
				return err
			}
			// `Aztec` verification
			if roleDax.Id == 1 {
				unixNow := time.Now().Unix()
				user.CreatedAt = &unixNow
			}
			// Role IDs
			currentRoleIds += fmt.Sprintf("%d,", roleDax.Id)
			// Circle
			roleIds = append(roleIds, roleDax.Id)
		}

		user.CurrentRoleIds = currentRoleIds

		var hasInnerCircleId bool = false
		var maxInnerOrderId int = -1
		for _, roleId := range roleIds {
			circle, order := getCircleAndOrderFromRoleId(roleId)
			fmt.Println(circle, order)
			if circle == 1 {
				hasInnerCircleId = true
				if order > maxInnerOrderId {
					maxInnerOrderId = order
				}
			}
		}

		if hasInnerCircleId {
			user.CurrentCircle = "INNER"
		} else {
			user.CurrentCircle = "OUTER"
		}

		if maxInnerOrderId == -1 {
			user.CurrentInnerOrder = nil
		} else {
			user.CurrentInnerOrder = &maxInnerOrderId
		}

		updatedUser, updateErr := usersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error udpating user in DB:", err)
			return err
		}
		fmt.Printf("User with CurrentRoleIds: %s\n", updatedUser.CurrentRoleIds)

		return nil
	}

	return nil
}

func getCircleAndOrderFromRoleId(roleId int) (int, int) {

	if roleId <= 7 {
		return 0, -1
	} else {
		if roleId >= 7 && roleId < 12 {
			return 1, 1
		} else if roleId >= 12 && roleId < 15 {
			return 1, 2
		} else if roleId >= 15 {
			return 1, 3
		}
	}

	return 0, -1

}
