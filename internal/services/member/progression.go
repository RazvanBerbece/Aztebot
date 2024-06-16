package member

import (
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
)

/*
|-------------------|--------------|----------------|-------------------|
|      Order        |      XP      | HOURS ON VOICE | MESSAGES SENT     |
|-------------------|--------------|----------------|-------------------|
| First Order       												    |
|-------------------|--------------|----------------|-------------------|
| Zelator           | 7,500 XP     | 15 HOURS       | 1,000 MESSAGES    |
|-------------------|--------------|----------------|-------------------|
| Theoricus         | 10,000 XP    | 20 HOURS       | 2,500 MESSAGES    |
|-------------------|--------------|----------------|-------------------|
| Philosophus       | 15,000 XP    | 30 HOURS       | 5,000 MESSAGES    |
|-------------------|--------------|----------------|-------------------|
| Second Order      												    |
|-------------------|--------------|----------------|-------------------|
| Adeptus Minor     | 30,000 XP    | 40 HOURS       | 12,500 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
| Adeptus Major     | 45,000 XP    | 45 HOURS       | 15,000 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
| Adeptus Exemptus  | 50,000 XP    | 50 HOURS       | 20,000 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
| Third Order       | 												    |
|-------------------|--------------|----------------|-------------------|
| Magister Templi   | 100,000 XP   | 200 HOURS      | 35,000 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
| Magus             | 150,000 XP   | 250 HOURS      | 45,000 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
| Ipsississimus     | 200,000 XP   | 300 HOURS      | 50,000 MESSAGES   |
|-------------------|--------------|----------------|-------------------|
*/

// Checks whether for the current stats of a member (XP, messages, time spent in VCs, etc.)
// a progression is due and sends the associated event to kickstart the promotion for the member.
func ProcessProgressionForMember(userId string, guildId string) error {

	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		return err
	}

	// Check if for the given member stats, a promotion is in order
	stats, err := globalRepositories.UserStatsRepository.GetStatsForUser(userId)
	if err != nil {
		return err
	}
	// Frame the given stats into one of the pre-defined ranks
	var xp = user.CurrentExperience
	var currentLevel = user.CurrentLevel
	var messagesSent = stats.NumberMessagesSent
	var timeSpentInVoice = stats.TimeSpentInVoiceChannels

	// Send event to process supposed promotion for user with updated stats
	globalMessaging.PromotionRequestsChannel <- events.PromotionRequestEvent{
		GuildId:       guildId,
		UserId:        userId,
		UserTag:       user.DiscordTag,
		CurrentLevel:  currentLevel,
		CurrentXp:     xp,
		MessagesSent:  messagesSent,
		TimeSpentInVc: timeSpentInVoice,
	}

	return nil
}

func GetRoleNameAndLevelFromStats(userXp float64, userNumberMessagesSent int, userTimeSpentInVc int) (string, int) {

	const sHour = 60 * 60

	// Check current stats against progression table
	// Figure out the promoted role to be given
	var level int = 0
	var roleName string = ""
	switch {
	// No order
	case userXp < 7500:
		// outer circle, so no role or level
	// First order
	case userXp >= 7500 && userXp < 10000:
		if userNumberMessagesSent >= 1000 && userTimeSpentInVc >= sHour*15 {
			level = 1
			roleName = "🔗 Zelator"
		}
	case userXp >= 10000 && userXp < 15000:
		if userNumberMessagesSent >= 2500 && userTimeSpentInVc >= sHour*20 {
			level = 2
			roleName = "📖 Theoricus"
		}
	case userXp >= 15000 && userXp < 30000:
		if userNumberMessagesSent >= 5000 && userTimeSpentInVc >= sHour*30 {
			level = 3
			roleName = "📿 Philosophus"
		}
	// Second order
	case userXp >= 30000 && userXp < 45000:
		if userNumberMessagesSent >= 12500 && userTimeSpentInVc >= sHour*40 {
			level = 4
			roleName = "🔮 Adeptus Minor"
		}
	case userXp >= 45000 && userXp < 50000:
		if userNumberMessagesSent >= 15000 && userTimeSpentInVc >= sHour*45 {
			level = 5
			roleName = "〽️ Adeptus Major"
		}
	case userXp >= 50000 && userXp < 100000:
		if userNumberMessagesSent >= 20000 && userTimeSpentInVc >= sHour*50 {
			level = 6
			roleName = "🧿 Adeptus Exemptus"
		}
	// Third order
	case userXp >= 100000 && userXp < 150000:
		if userNumberMessagesSent >= 35000 && userTimeSpentInVc >= sHour*200 {
			level = 7
			roleName = "☀️ Magister Templi"
		}
	case userXp >= 150000 && userXp < 200000:
		if userNumberMessagesSent >= 45000 && userTimeSpentInVc >= sHour*250 {
			level = 8
			roleName = "🧙🏼 Magus"
		}
	case userXp >= 200000:
		if userNumberMessagesSent >= 50000 && userTimeSpentInVc >= sHour*300 {
			level = 9
			roleName = "⚔️ Ipsissimus"
		}
	}

	return roleName, level
}
