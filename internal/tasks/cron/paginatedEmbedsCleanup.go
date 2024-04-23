package cron

import (
	"fmt"
	"time"

	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	actionEventsUtils "github.com/RazvanBerbece/Aztebot/internal/handlers/remoteEvents/actionEvents/utils"
	"github.com/bwmarrin/discordgo"
)

func ClearOldPaginatedEmbeds(s *discordgo.Session) {

	var numSec int = 20 // 60 * 15

	fmt.Println("[CRON] Starting Task ClearOldPaginatedEmbeds() at", time.Now(), "running every", numSec, "seconds")

	ticker := time.NewTicker(time.Duration(numSec) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go cleanupOldPaginatedEmbeds(s)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func cleanupOldPaginatedEmbeds(s *discordgo.Session) {

	for msgId, embedData := range globalState.EmbedsToPaginate {
		// If old enough
		threshold := time.Second * 5 // 10 minutes ?
		if time.Since(time.Unix(int64(embedData.Timestamp), 0)) > threshold {
			// Remove action row from embed
			go actionEventsUtils.DisablePaginatedEmbed(s, embedData.ChannelId, msgId)
			// Remove from global map
			delete(globalState.EmbedsToPaginate, msgId)
		}
	}

}
