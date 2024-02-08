package readyEvent

import (
	"fmt"
	"strings"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/globals"
)

func LoadStaticData() {
	LoadNotificationChannels()
}

// Load the available notification channels in the global map.
func LoadNotificationChannels() {

	for _, channelPairString := range globals.NotificationChannelsPairs {

		isVoice, descriptor, channelId := getChannelValuesFromChannelPair(channelPairString)

		if descriptor == nil || channelId == nil {
			continue
		}

		globals.NotificationChannels[*descriptor] = dataModels.Channel{
			IsVoice:    isVoice,
			Descriptor: *descriptor,
			ChannelId:  *channelId,
		}

	}

}

// Returns a channel details as 3 variables (IsVoiceChannel, Descriptor, ChannelId) as parsed from the input
// represented by a string of format [ "descriptor channelId" ] (e.g. "notif-timeout 1234567890")
func getChannelValuesFromChannelPair(channelPair string) (bool, *string, *string) {

	var channelValues = strings.Split(channelPair, " ")

	fmt.Println(channelValues)

	if len(channelValues) != 2 {
		return false, nil, nil
	}

	fmt.Println(channelValues)

	// Figure out if channel might be a voice channel
	var isVoice = false
	// TODO

	return isVoice, &channelValues[0], &channelValues[1]

}
