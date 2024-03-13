package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ProcessCleanupUnusedDynamicChannels(s *discordgo.Session, guildId string) {

	var numSec int = 3 // 60 * 60 * 5

	fmt.Println("[CRON] Starting Cron Ticker CleanupUnusedDynamicChannels() at", time.Now(), "running every", numSec, "seconds")

	// Run on startup too, it can't hurt :)
	CleanupUnusedDynamicChannels(s, guildId)

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				CleanupUnusedDynamicChannels(s, guildId)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func CleanupUnusedDynamicChannels(s *discordgo.Session, guildId string) {
	channels, err := s.GuildChannels(guildId)
	if err != nil {
		fmt.Printf("An error ocurred while cleaning up hanging dynamic channels: %v\n", err)
		return
	}

	for _, channel := range channels {
		// If channel is a dynamic channel - given the `(~Extra~)` substring in the name
		if channel.Type == discordgo.ChannelTypeGuildVoice && strings.Contains(channel.Name, "(~Extra~)") {
			// If the channel is empty
			if channel.MemberCount == 0 {
				// Then delete it
				_, err := s.ChannelDelete(channel.ID)
				if err != nil {
					fmt.Printf("An error ocurred while cleaning up hanging dynamic channels: %v\n", err)
					return
				}
			}
		}
	}
}
