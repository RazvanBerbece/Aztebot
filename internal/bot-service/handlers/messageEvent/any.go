package messageEvent

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/bot-service/api/member"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/bot-service/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func Any(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageCreatorUserId := m.Author.ID

	// Ignore all messages created by bots
	authorIsBot, err := member.MemberIsBot(s, globals.DiscordMainGuildId, messageCreatorUserId)
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

	// Grant experience points
	go member.GrantMemberExperience(messageCreatorUserId, "MSG_REWARD", nil)

}
