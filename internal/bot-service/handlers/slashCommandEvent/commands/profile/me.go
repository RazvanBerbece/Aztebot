package profileSlashHandlers

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	rolesService "github.com/RazvanBerbece/Aztebot/internal/bot-service/api/roles"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMe(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userId := i.Interaction.Member.User.ID

	// Attempt a sync
	err := ProcessUserUpdate(userId, s, i)
	if err != nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, fmt.Sprintf("An error ocurred while trying to sync your profile card: `%s`", err))
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: utils.SimpleEmbed("🤖   Slash Command Confirmation", "Gathering your `/me` data..."),
		},
	})

	embed := GetProfileEmbedForUser(s, userId)
	if embed == nil {
		utils.ErrorEmbedResponseEdit(s, i.Interaction, "An error ocurred while trying to fetch your profile card.")
		return
	}

	// Final response
	editContent := ""
	editWebhook := discordgo.WebhookEdit{
		Content: &editContent,
		Embeds:  &embed,
	}
	s.InteractionResponseEdit(i.Interaction, &editWebhook)
}

func GetProfileEmbedForUser(s *discordgo.Session, userId string) []*discordgo.MessageEmbed {

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

	// Process highest roles
	var highestOrderRole *dataModels.Role = nil
	var highestStaffRole *dataModels.Role = nil
	var orderRoleText string = ""
	roles, err := globalsRepo.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		highestOrderRole = nil
		highestStaffRole = nil
	} else {
		highestStaffRole, highestOrderRole = rolesService.GetHighestRoles(roles)
		orderRoleText = fmt.Sprintf("%s | ", highestOrderRole.DisplayName)
	}

	stats, errStats := globalsRepo.UserStatsRepository.GetStatsForUser(userId)
	if errStats != nil {
		log.Fatalf("Cannot retrieve user %s stats: %v", user.DiscordTag, errStats)
		return nil
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

	var genderDisplayString = ""
	switch user.Gender {
	case 0:
		genderDisplayString = " ♂"
	case 1:
		genderDisplayString = " ♀"
	case 2:
		genderDisplayString = " ⚥"
	case 3:
		genderDisplayString = " 🌈"
	default:
		// undefined gender - not friendly for displaying purposes
	}

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("🤖   `%s`'s Profile Card%s", user.DiscordTag, genderDisplayString)).
		SetDescription(fmt.Sprintf("`%s%s CIRCLE%s`", orderRoleText, user.CurrentCircle, orderText)).
		SetThumbnail(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userId, apiUser.Avatar)).
		SetColor(000000).
		AddLineBreakField()

	if userCreatedTimeString != "" && highestOrderRole != nil {

		// Process ranks in leaderboards
		msgRankString := ""
		reactRankString := ""
		streakRankString := ""
		vcRankString := ""
		musicRankString := ""
		xpRankString := ""
		ranks, err := member.GetMemberRankInLeaderboards(userId)
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

		// Retrieve experience points for user
		xp, err := member.GetMemberExperiencePoints(userId)
		if err != nil {
			log.Printf("Cannot retrieve user %s XP: %v", userId, err)
			return nil
		}
		xpInt := int(*xp)
		xpRank, rankErr := member.GetMemberXpRank(userId)
		if rankErr != nil {
			log.Printf("Cannot retrieve user %s XP rank: %v", userId, rankErr)
			return nil
		}
		xpRankString = fmt.Sprintf(" (`🏆 #%d`)", *xpRank)

		// Add extra decorations to the embed (special users, staff members, etc.)
		DecorateProfileEmbed(embed, highestStaffRole, userId)

		embed.
			AddField(fmt.Sprintf("🩸 Aztec since:  `%s`", userCreatedTimeString), "", false).
			AddField(fmt.Sprintf("🔄 Active day streak:  `%d`%s", stats.NumberActiveDayStreak, streakRankString), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("✉️ Total messages sent:  `%d`%s", stats.NumberMessagesSent, msgRankString), "", false).
			AddField(fmt.Sprintf("⚙️ Total slash commands used:  `%d`", stats.NumberSlashCommandsUsed), "", false).
			AddField(fmt.Sprintf("💯 Total reactions received:  `%d`%s", stats.NumberReactionsReceived, reactRankString), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("🎙️ Time spent in voice channels:  `%s`%s", timeSpentInVcs, vcRankString), "", false).
			AddField(fmt.Sprintf("🎵 Time spent listening to music:  `%s`%s", timeSpentListeningMusic, musicRankString), "", false).
			AddLineBreakField().
			AddField(fmt.Sprintf("💠 Total gained XP:  `%d`%s", xpInt, xpRankString), "", false)

	} else {
		embed.AddField("Member hasn't verified yet.", "", false)
	}

	return []*discordgo.MessageEmbed{embed.MessageEmbed}
}

func DecorateProfileEmbed(embed *embed.Embed, staffRole *dataModels.Role, userId string) {

	// Special users segment
	if userId == "526512064794066945" {
		// The one and only, Edi
		embed.AddField("👑 Azteca", "", false)
	}

	// Staff text segment (is user a member of staff?) in embed description
	if member.IsStaff(userId) {
		var staffFieldName string = "💎 OTA Staff Member"
		if staffRole != nil {
			staffFieldName += fmt.Sprintf(" (`%s`)", staffRole.DisplayName)
		}
		embed.AddField(staffFieldName, "", false)
	}

}
