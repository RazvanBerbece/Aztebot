package messageEvent

import (
	"fmt"
	"strings"

	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/bwmarrin/discordgo"
)

func SimpleListenReply(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.ToLower(m.Content) == "aztebot" {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
	}

	remoteTag := m.Content
	localTag := fmt.Sprintf("<@%s>", globalConfiguration.DiscordAztebotAppId)
	if remoteTag == localTag {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
	}

}
