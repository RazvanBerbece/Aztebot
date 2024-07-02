package profileSlashHandlers

import (
	"fmt"
	"log"
	"time"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	rolesService "github.com/RazvanBerbece/Aztebot/internal/services/roles"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleSlashMe(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userId := i.Interaction.Member.User.ID

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

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user with id %s: %v\n", userId, err)
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
	var highestOrderRole *dax.Role = nil
	var highestStaffRole *dax.Role = nil
	var orderRoleText string = ""
	roles, err := globalRepositories.UsersRepository.GetRolesForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve roles for user with id %s: %v", userId, err)
		highestOrderRole = nil
		highestStaffRole = nil
	} else {
		highestStaffRole, highestOrderRole = rolesService.GetHighestRoles(roles)
		if highestOrderRole != nil {
			orderRoleText = fmt.Sprintf("%s | ", highestOrderRole.DisplayName)
		}
	}

	stats, errStats := globalRepositories.UserStatsRepository.GetStatsForUser(userId)
	if errStats != nil {
		log.Fatalf("Cannot retrieve user %s stats: %v", user.DiscordTag, errStats)
		return nil
	}

	// Process the time spent in VCs in a nice format
	sTimeSpentInVc := int64(stats.TimeSpentInVoiceChannels)
	daysVC, hoursVC, minutesVC, secondsVC := utils.HumanReadableDuration(float64(sTimeSpentInVc))
	timeSpentInVcs := fmt.Sprintf("%dd, %dh:%dm:%ds", daysVC, hoursVC, minutesVC, secondsVC)

	// Process the time spent listening to music a nice format
	sTimeSpentListeningMusic := int64(stats.TimeSpentListeningToMusic)
	daysMusic, hoursMusic, minutesMusic, secondsMusic := utils.HumanReadableDuration(float64(sTimeSpentListeningMusic))
	timeSpentListeningMusic := fmt.Sprintf("%dd, %dh:%dm:%ds", daysMusic, hoursMusic, minutesMusic, secondsMusic)

	// Get the profile picture url
	// Fetch user information from Discord API.
	apiUser, err := s.User(userId)
	if err != nil {
		log.Printf("Cannot retrieve user %s from Discord API: %v\n", userId, err)
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

	var repText string = ""
	userRep, err := globalRepositories.UserRepRepository.GetRepForUser(userId)
	if err != nil {
		log.Printf("Cannot retrieve user rep %s from DB: %v\n", userId, err)
		return nil
	}
	if userRep != nil {
		if userRep.Rep > 0 {
			repText = fmt.Sprintf("\n+%d rep", userRep.Rep)
		} else if userRep.Rep < 0 {
			repText = fmt.Sprintf("\n%d rep", userRep.Rep)
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
		SetDescription(fmt.Sprintf("`%s%s CIRCLE%s%s`", orderRoleText, user.CurrentCircle, orderText, repText)).
		SetThumbnail(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userId, apiUser.Avatar)).
		SetColor(000000).
		DecorateWithTimestampFooter("Mon, 02 Jan 2006 15:04:05 MST").
		AddLineBreakField()

	if member.IsFullyVerifiedLocal(*user) {

		// Process ranks in leaderboards
		msgRankString := ""
		reactRankString := ""
		streakRankString := ""
		vcRankString := ""
		musicRankString := ""
		xpRankString := ""
		ranks, err := member.GetMemberRankInLeaderboards(userId)
		if err != nil {
			log.Printf("Cannot retrieve user %s leaderboard ranks: %v\n", userId, err)
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

		// XP value and XP rank
		xpInt := int(user.CurrentExperience)
		xpRank, rankErr := member.GetMemberXpRank(userId)
		if rankErr != nil {
			log.Printf("Cannot retrieve user %s XP rank: %v\n", userId, rankErr)
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

func DecorateProfileEmbed(embed *embed.Embed, staffRole *dax.Role, userId string) {

	// Special users segment
	if userId == "526512064794066945" {
		// The one and only, Edi
		embed.AddField("👑 Azteca", "", false)
	}

	// Staff text segment (is user a member of staff?) in embed description
	if member.IsStaff(userId, globalConfiguration.StaffRoles) {
		var staffFieldName string = "💎 OTA Staff Member"
		if staffRole != nil {
			staffFieldName += fmt.Sprintf(" (`%s`)", staffRole.DisplayName)
		}
		embed.AddField(staffFieldName, "", false)
	}

}
