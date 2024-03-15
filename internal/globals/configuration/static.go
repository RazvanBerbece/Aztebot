package globalConfiguration

import dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"

var JailTasks = []string{}

var DefaultInformationChannels map[string]string
var AfkChannels map[string]string
var MusicChannels map[string]string
var DynamicChannelCreateButtonIds map[string]string

var NotificationChannels = make(map[string]dataModels.Channel)
