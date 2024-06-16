package slashHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
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

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
		if err == sql.ErrNoRows {
			log.Printf("Storing user with id %s", userId)
			user, err = globalsRepo.UsersRepository.SaveInitialUserDetails(event.Member.User.Username, userId)
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
			roleDax, err := globalsRepo.RolesRepository.GetRole(userRoleObj.Name)
			if err != nil {
				if err == sql.ErrNoRows {
					// This will probably be a role which is assigned to the three orders or something, so we can ignore
					// and move on to the other roles of the user
					continue
				} else {
					log.Printf("Error getting role %s from DB: %v", userRoleObj.Name, err)
					return err
				}
			} else {
				// `Aztec` verification
				if roleDax.Id == 1 && user.CreatedAt == nil {
					unixNow := time.Now().Unix()
					user.CreatedAt = &unixNow
				}
				// Role IDs
				currentRoleIds += fmt.Sprintf("%d,", roleDax.Id)
				// Circle
				roleIds = append(roleIds, roleDax.Id)
			}
		}

		user.CurrentRoleIds = currentRoleIds

		var hasInnerCircleId bool = false
		var maxInnerOrderId int = -1
		for _, roleId := range roleIds {
			circle, order := utils.GetCircleAndOrderFromRoleId(roleId)
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

		_, updateErr := globalsRepo.UsersRepository.UpdateUser(*user)
		if updateErr != nil {
			log.Println("Error updating user in DB:", err)
			return err
		}

		return nil
	}

	return nil
}
