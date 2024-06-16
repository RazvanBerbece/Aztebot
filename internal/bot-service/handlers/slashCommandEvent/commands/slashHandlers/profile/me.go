package profile

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Gathering your `/me` data..."),
		},
	})

	embed := displayEmbedForUser(s, i.Interaction.Member.User.ID)
	if embed == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while trying to fetch your profile card.",
			},
		})
	}

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &embed,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)
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
	stats, errStats := globalsRepo.UserStatsRepository.GetStatsForUser(userId)
	if errStats != nil {
		if errStats == sql.ErrNoRows {
			errStatsInit := globalsRepo.UserStatsRepository.SaveInitialUserStats(userId)
			if errStatsInit != nil {
				log.Fatalf("Cannot store initial user %s stats: %v", user.DiscordTag, errStatsInit)
				return nil
			}
		}
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

	// Process the time spent listening to music a nice format
	sTimeSpentListeningMusic := int64(stats.TimeSpentListeningToMusic)
	daysMusic, hoursMusic, minutesMusic, secondsMusic := utils.HumanReadableTimeLength(float64(sTimeSpentListeningMusic))
	timeSpentListeningMusic := fmt.Sprintf("%dd, %dh:%dm:%ds", daysMusic, hoursMusic, minutesMusic, secondsMusic)

	// Get the profile picture url
	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		log.Printf("Cannot retrieve user %s from Discord API: %v", userId, err)
		return nil
	}

	var orderText string = ""
	if user.CurrentInnerOrder != nil {
		orderNum := *user.CurrentInnerOrder
		switch orderNum {
		case 1:
			orderText = " | FIRST ORDER"
		case 2:
			orderText = " | SECOND ORDER"
		case 3:
			orderText = " | THIRD ORDER"
		}
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Profile Card", user.DiscordTag)).
		SetDescription(fmt.Sprintf("`%s CIRCLE%s`", user.CurrentCircle, orderText)).
		SetThumbnail(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userId, apiUser.Avatar)).
		SetColor(000000).
		AddLineBreakField()

	if userCreatedTimeString != "" && highestRole != nil {

		// Process ranks in leaderboards
		msgRankString := ""
		reactRankString := ""
		streakRankString := ""
		vcRankString := ""
		musicRankString := ""
		ranks, err := member.GetMemberRankInLeaderboards(s, userId)
		if err != nil {
			log.Printf("Cannot retrieve user %s leaderboard ranks: %v", userId, err)
			return nil
		}
		if msgRank, ok := ranks["msg"]; ok {
			msgRankString = fmt.Sprintf(" (`🏆 #%d`)", msgRank)
		}
		if reactRank, ok := ranks["react"]; ok {
			reactRankString = fmt.Sprintf(" (`🏆 #%d`)", reactRank)
		}
		if streakRank, ok := ranks["streak"]; ok {
			streakRankString = fmt.Sprintf(" (`🏆 #%d`)", streakRank)
		}
		if vcRank, ok := ranks["vc"]; ok {
			vcRankString = fmt.Sprintf(" (`🏆 #%d`)", vcRank)
		}
		if musicRank, ok := ranks["music"]; ok {
			musicRankString = fmt.Sprintf(" (`🏆 #%d`)", musicRank)
		}

		if user.UserId == "526512064794066945" {
			// The one and only, Edi
			embed.AddField("👑 Azteca", "", false)
		}

		if isStaffMember {
			embed.AddField("💎 OTA Staff Member", "", false)
		}

		embed.
			AddField(fmt.Sprintf("🩸 Aztec since:  `%s`", userCreatedTimeString), "", false).
			AddField(fmt.Sprintf("⭐ Highest obtained role:  `%s`", highestRole.DisplayName), "", false).
			AddField(fmt.Sprintf("🔄 Active day streak:  `%d`%s", stats.NumberActiveDayStreak, streakRankString), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("✉️ Total messages sent:  `%d`%s", stats.NumberMessagesSent, msgRankString), "", false).
			AddField(fmt.Sprintf("⚙️ Total slash commands used:  `%d`", stats.NumberSlashCommandsUsed), "", false).
			AddField(fmt.Sprintf("💯 Total reactions received:  `%d`%s", stats.NumberReactionsReceived, reactRankString), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("🎙️ Time spent in voice channels:  `%s`%s", timeSpentInVcs, vcRankString), "", false).
			AddField(fmt.Sprintf("🎵 Time spent listening to music:  `%s`%s", timeSpentListeningMusic, musicRankString), "", false)

	} else {
		embed.AddField("Member hasn't verified yet.", "", false)
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}
