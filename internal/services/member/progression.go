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
| Practicus         | 15,000 XP    | 25 HOURS       | 3,500 MESSAGES    |
|-------------------|--------------|----------------|-------------------|
| Philosophus       | 20,000 XP    | 30 HOURS       | 5,000 MESSAGES    |
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
	var messagesSent = stats.NumberMessagesSent
	var timeSpentInVoice = stats.TimeSpentInVoiceChannels

	// Send event to process supposed promotion for user with updated stats
	globalMessaging.PromotionRequestsChannel <- events.PromotionRequestEvent{
		GuildId:       guildId,
		UserId:        userId,
		CurrentXp:     xp,
		MessagesSent:  messagesSent,
		TimeSpentInVc: timeSpentInVoice,
	}

	return nil
}
