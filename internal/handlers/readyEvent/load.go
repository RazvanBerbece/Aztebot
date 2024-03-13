package readyEvent

import (
	"strings"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
)

func LoadStaticData() {
	LoadNotificationChannels()
	LoadJailTasks()
	LoadStaticDiscordChannels()
}

// Load the available tasks to get out of Jail in the global list.
func LoadJailTasks() {
	globals.JailTasks = []string{
		"Continue the lyrics",
		"Write a Poem",
		"Math Quiz",
		"AzteQuiz",
		"Guess the Number",
		"Roll a Double",
	}
}

// Load some static Discord channel IDs (useful for main guild ops)
func LoadStaticDiscordChannels() {

	if globals.Environment == "staging" {
		// Dev afk channels
		globals.AfkChannels = map[string]string{
			"1176284686297874522": "afk",
		}
	} else {
		// Production afk channels
		globals.AfkChannels = map[string]string{
			"1212508073101627412": "afk",
		}
	}

	if globals.Environment == "staging" {
		// Dev music channels
		globals.MusicChannels = map[string]string{
			"1173790229258326106": "radio",
		}
	} else {
		// Production music channels
		globals.MusicChannels = map[string]string{
			"1176204022399631381": "radio",
			"1118202946455351388": "music-1",
			"1118202975026937948": "music-2",
			"1118202999504904212": "music-3",
		}
	}

	if globals.Environment == "staging" {
		// Dev dynamic channel creation button channels
		globals.DynamicChannelCreateButtonIds = map[string]string{
			"1217251206624186481": "☕ | Dev Test Room",
		}
	} else {
		// Production dynamic channel creation button channels
		globals.DynamicChannelCreateButtonIds = map[string]string{
			"1171570400891785266": "☕ | Chill Room (~Extra~)",
			"1171589545473613886": "🔒 | Private Room (~Extra~)",
			"1171591013354197062": "🔮 | Spiritual Room (~Extra~)",
			"1171595498185035796": "🎵 | Music Room (~Extra~)",
			"1171599680568832023": "🎮 | Gaming (~Extra~)",
		}
	}

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

	if len(channelValues) != 2 {
		return false, nil, nil
	}

	// Figure out if channel might be a voice channel
	var isVoice = false
	// TODO

	return isVoice, &channelValues[0], &channelValues[1]

}
