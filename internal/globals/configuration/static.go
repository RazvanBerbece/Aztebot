package globalConfiguration

import "github.com/RazvanBerbece/Aztebot/internal/data/models/domain"

var JailTasks = []string{}

var DefaultInformationChannels map[string]string
var AfkChannels map[string]string
var MusicChannels map[string]string
var DynamicChannelCreateButtonIds map[string]string

var DeleteExceptedChannels map[string]string

var NotificationChannels = make(map[string]domain.Channel)
