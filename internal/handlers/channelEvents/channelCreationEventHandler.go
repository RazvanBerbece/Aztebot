package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	server_channel "github.com/RazvanBerbece/Aztebot/internal/services/server/channel"
	"github.com/bwmarrin/discordgo"
)

func HandleDynamicChannelCreationEvents(s *discordgo.Session) {

	for channelEvent := range globalMessaging.ChannelCreationsChannel {

		// Limit dynamic channels to a maximum of ~25, to minimise the risk of DoS by spamming VCs
		if globalState.DynamicChannelsCount >= 25 {
			continue
		}

		categoryId, err := server_channel.GetCategoryIdForChannel(s, channelEvent.ParentGuildId, channelEvent.ParentChannelId)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (get parent category): %v\n", err)
			continue
		}

		// Create a new voice channel with the given specification
		createdChannel, err := server_channel.CreateVoiceChannelForCategory(s, channelEvent.ParentGuildId, categoryId, channelEvent.Name, channelEvent.Private)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (create channel): %v\n", err)
			continue
		}

		globalState.DynamicChannelsCount += 1

		// and move member to it
		err = s.GuildMemberMove(channelEvent.ParentGuildId, channelEvent.ParentMemberId, &createdChannel.ID)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (move member): %v\n", err)
			continue
		}
	}

}
