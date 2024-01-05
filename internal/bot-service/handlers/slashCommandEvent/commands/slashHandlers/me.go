package slashHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMe(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Attempt a sync
	err := ProcessUserUpdate(i.Interaction.Member.User.ID, s, i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your profile card.",
			},
		})
	}

	embed := displayEmbedForUser(s, i.Interaction.Member.User.ID)
	if embed == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your profile card.",
			},
		})
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})
}

func displayEmbedForUser(s *discordgo.Session, userId string) []*discordgo.MessageEmbed {

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v", userId, err)
		return nil
	}

	// Format CreatedAt
	var userCreatedTime time.Time
	var userCreatedTimeString string
	if user.CreatedAt != nil {
		userCreatedTime = time.Unix(*user.CreatedAt, 0).UTC()
		userCreatedTimeString = userCreatedTime.Format("January 2, 2006")
	} else {
		userCreatedTimeString = ""
	}

	// Process highest role
	var highestRole *dataModels.Role
	roles, err := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		highestRole = nil
	} else {
		highestRole = &roles[len(roles)-1] // role IDs for users are stored in DB in ascending order by rank, so the last one is the highest
	}

	// Setup user stats if the user doesn't have an entity in UserStats
	_, errStats := globalsRepo.UserStatsRepository.GetStatsForUser(userId)
	if errStats != nil {
		if errStats == sql.ErrNoRows {
			errStatsInit := globalsRepo.UserStatsRepository.SaveInitialUserStats(userId)
			if errStatsInit != nil {
				log.Fatalf("Cannot store initial user %s stats: %v", user.DiscordTag, errStatsInit)
				return nil
			}
		}
	}

	// Get user stats
	stats, err := globalsRepo.UserStatsRepository.GetStatsForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user %s stats from DB: %v", userId, err)
		return nil
	}

	// Staff text segment (is user a member of staff?) in embed description
	var isStaffMember bool = false
	for _, role := range roles {
		if role.Id == 3 || role.Id == 5 || role.Id == 6 || role.Id == 7 || role.Id == 18 {
			// User is a staff member if they belong to any of the roles above
			isStaffMember = true
		}
	}

	// Process the time spent in VCs in a nice format
	sTimeSpentInVc := int64(stats.TimeSpentInVoiceChannels)
	daysVC, hoursVC, minutesVC, secondsVC := utils.HumanReadableTimeLength(float64(sTimeSpentInVc))
	timeSpentInVcs := fmt.Sprintf("%dd, %dh:%dm:%ds", daysVC, hoursVC, minutesVC, secondsVC)

	// Get the profile picture url
	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		return nil
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Profile Card", user.DiscordTag)).
		SetDescription(fmt.Sprintf("`%s CIRCLE`", user.CurrentCircle)).
		SetThumbnail(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userId, apiUser.Avatar)).
		SetColor(000000)

	if userCreatedTimeString != "" && highestRole != nil {

		if isStaffMember {
			embed.AddField("💎 OTA Staff Member", "", false)
		}

		embed.
			AddField(fmt.Sprintf("🩸 Aztec since:  `%s`", userCreatedTimeString), "", false).
			AddField(fmt.Sprintf("⭐ Highest obtained role:  `%s`", highestRole.DisplayName), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("✉️ Total messages sent:  `%d`", stats.NumberMessagesSent), "", false).
			AddField(fmt.Sprintf("⚙️ Total slash commands used:  `%d`", stats.NumberSlashCommandsUsed), "", false).
			AddField(fmt.Sprintf("💯 Total reactions received:  `%d`", stats.NumberReactionsReceived), "", false).
			AddField(fmt.Sprintf("🔄 Active day streak:  `%d`", stats.NumberActiveDayStreak), "", false).
			AddField(fmt.Sprintf("🎙️ Time spent in voice channels:  `%s`", timeSpentInVcs), "", false).
			AddLineBreakField().
			AddField("", "_(Stats collected after 15/12/2023)_", false)

	} else {
		embed.AddField("Member hasn't verified yet.", "", false)
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
