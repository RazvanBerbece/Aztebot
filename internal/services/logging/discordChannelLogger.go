package logging

import (
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	"github.com/bwmarrin/discordgo"
)

type DiscordChannelLogger struct {
	Session *discordgo.Session
	Topic   string
}

func NewDiscordLogger(s *discordgo.Session, topic string) *DiscordChannelLogger {
	return &DiscordChannelLogger{
		Session: s,
		Topic:   topic,
	}
}

func (l DiscordChannelLogger) LogInfo(msg string) {
	if channel, channelExists := globalConfiguration.NotificationChannels[l.Topic]; channelExists {
		globalMessaging.NotificationsChannel <- events.NotificationEvent{
			TargetChannelId: channel.ChannelId,
			Type:            "DEFAULT",
			TextData:        &msg,
		}
	}
}
