package channelHandlers

import (
	"fmt"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandleDirectMessageEvents(s *discordgo.Session) {
	for directMessageEvent := range globalMessaging.DirectMessagesChannel {
		if directMessageEvent.Embed != nil {
			// The event has an embed to passthrough
			if directMessageEvent.PaginationRow != nil {
				// and the embed supports pagination !
				// so add the action row to the request
				err := member.SendDirectComplexEmbedToMember(s, directMessageEvent.UserId, *directMessageEvent.Embed, *directMessageEvent.PaginationRow, globalConfiguration.EmbedPageSize)
				if err != nil {
					fmt.Printf("Failed to process DirectMessageEvent (Pagination: On): %v\n", err)
				}
			} else {
				err := member.SendDirectEmbedToMember(s, directMessageEvent.UserId, *directMessageEvent.Embed)
				if err != nil {
					fmt.Printf("Failed to process DirectMessageEvent: %v\n", err)
				}
			}
		} else {
			if directMessageEvent.Text != nil && directMessageEvent.Title != nil {
				// The event has a title and text to use in a dynamic embed
				err := member.SendDirectSimpleEmbedToMember(s, directMessageEvent.UserId, *directMessageEvent.Title, *directMessageEvent.Text)
				if err != nil {
					fmt.Printf("Failed to process DirectMessageEvent: %v\n", err)
				}
			} else {
				fmt.Println("This DM event:", directMessageEvent, "is not valid.")
			}
		}
	}
}
