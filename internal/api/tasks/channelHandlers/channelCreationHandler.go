package channelHandlers

import (
	"fmt"

	"github.com/RazvanBerbece/Aztebot/internal/api/server"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	"github.com/bwmarrin/discordgo"
)

func HandleChannelCreationMessages(s *discordgo.Session) {

	for channelEvent := range globals.ChannelCreationsChannel {

		// Limit dynamic channels to a maximum of ~25, to minimise the risk of DoS by spamming VCs
		if globals.DynamicChannelsCount >= 25 {
			continue
		}

		categoryId, err := server.GetCategoryIdForChannel(s, channelEvent.ParentGuildId, channelEvent.ParentChannelId)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (get parent category): %v\n", err)
			continue
		}

		// Create a new voice channel with the given specification
		createdChannel, err := server.CreateVoiceChannelForCategory(s, channelEvent.ParentGuildId, categoryId, channelEvent.Name, channelEvent.Private)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (create channel): %v\n", err)
			continue
		}

		globals.DynamicChannelsCount += 1

		// and move member to it
		err = s.GuildMemberMove(channelEvent.ParentGuildId, channelEvent.ParentMemberId, &createdChannel.ID)
		if err != nil {
			fmt.Printf("Failed to handle VC creation event (move member): %v\n", err)
			continue
		}
	}

}
