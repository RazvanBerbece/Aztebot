package globalMessaging

import "github.com/RazvanBerbece/Aztebot/internal/data/models/events"

var NotificationsChannel = make(chan events.NotificationEvent)
var ExperienceGrantsChannel = make(chan events.ExperienceGrantEvent)
var ChannelCreationsChannel = make(chan events.VoiceChannelCreateEvent)
var MessageDeletionChannel = make(chan events.MessageDeletionForUserEvent)
var DirectMessagesChannel = make(chan events.DirectMessageEvent)
var ComplexResponsesChannel = make(chan events.ComplexResponseEvent)
var PromotionRequestsChannel = make(chan events.PromotionRequestEvent)
