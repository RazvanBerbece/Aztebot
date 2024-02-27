package messageEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageCreatorUserId := m.Author.ID

	// Ignore all messages created by bots
	authorIsBot, err := member.IsBot(s, globals.DiscordMainGuildId, messageCreatorUserId, false)
	if err != nil {
		return
	}
	if *authorIsBot {
		return
	}

	// Increase stats for user
	err = globalsRepo.UserStatsRepository.IncrementMessagesSentForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) message count: %v\n", messageCreatorUserId, err)
	}

	err = globalsRepo.UserStatsRepository.IncrementActivitiesTodayForUser(messageCreatorUserId)
	if err != nil {
		fmt.Printf("An error ocurred while incrementing user (%s) activities count: %v\n", messageCreatorUserId, err)
	}
	err = globalsRepo.UserStatsRepository.UpdateLastActiveTimestamp(messageCreatorUserId, time.Now().Unix())
	if err != nil {
		fmt.Printf("An error ocurred while updating user (%s) last timestamp: %v\n", m.Author.ID, err)
	}

	// Publish experience grant message on the channel
	globals.ExperienceGrantsChannel <- dataModels.ExperienceGrant{
		UserId:   messageCreatorUserId,
		Points:   globals.ExperienceReward_MessageSent,
		Activity: "Message Sent Reward",
	}

}
