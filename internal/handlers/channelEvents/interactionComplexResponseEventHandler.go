package channelHandlers

import (
	"fmt"

	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/RazvanBerbece/Aztebot/internal/services/member"
	"github.com/bwmarrin/discordgo"
)

func HandleComplexResponseEvents(s *discordgo.Session, pageSize int) {
	for complexResponseEvent := range globalMessaging.ComplexResponsesChannel {
		if complexResponseEvent.Embed != nil {
			// The event has an embed to passthrough
			if len(complexResponseEvent.Embed.Fields) > pageSize && complexResponseEvent.PaginationRow != nil {
				// and the embed needs pagination !
				err := member.ReplyComplexToInteraction(s, complexResponseEvent.Interaction, *complexResponseEvent.Embed, *complexResponseEvent.PaginationRow, pageSize)
				if err != nil {
					fmt.Printf("Failed to process ComplexResponseEvent (Pagination: On): %v\n", err)
				}
			} else {
				editContent := ""
				editWebhook := discordgo.WebhookEdit{
					Content: &editContent,
					Embeds:  &[]*discordgo.MessageEmbed{complexResponseEvent.Embed.MessageEmbed},
				}
				s.InteractionResponseEdit(complexResponseEvent.Interaction, &editWebhook)
			}
		} else {
			if complexResponseEvent.Text != nil && complexResponseEvent.Title != nil {
				complexResponseEvent.Embed.Fields[0].Name = *complexResponseEvent.Title
				complexResponseEvent.Embed.Fields[0].Value = *complexResponseEvent.Text
				editContent := ""
				editWebhook := discordgo.WebhookEdit{
					Content: &editContent,
					Embeds:  &[]*discordgo.MessageEmbed{complexResponseEvent.Embed.MessageEmbed},
				}
				s.InteractionResponseEdit(complexResponseEvent.Interaction, &editWebhook)
			} else {
				fmt.Println("This response event:", complexResponseEvent, "is not valid.")
			}
		}
	}
}
