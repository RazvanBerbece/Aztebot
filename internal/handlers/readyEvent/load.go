package readyEvent

import (
	"strings"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
)

func LoadStaticData() {
	LoadNotificationChannels()
	LoadJailTasks()
	LoadStaticDiscordChannels()
}

// Load the available tasks to get out of Jail in the global list.
func LoadJailTasks() {
	globalConfiguration.JailTasks = []string{
		"Continue the lyrics",
		"Write a Poem",
		"Math Quiz",
		"AzteQuiz",
		"Guess the Number",
		"Roll a Double",
		"Write a Chorus for a Rap Song",
		"Custom Task from a Staff Member",
	}
}

// Load some static Discord channel IDs (useful for main guild ops)
func LoadStaticDiscordChannels() {

	if globalConfiguration.Environment == "staging" {
		// Dev afk channels
		globalConfiguration.AfkChannels = map[string]string{
			"1176284686297874522": "afk",
		}
	} else {
		// Production afk channels
		globalConfiguration.AfkChannels = map[string]string{
			"1212508073101627412": "afk",
		}
	}

	if globalConfiguration.Environment == "staging" {
		// Dev music channels
		globalConfiguration.MusicChannels = map[string]string{
			"1173790229258326106": "radio",
		}
	} else {
		// Production music channels
		globalConfiguration.MusicChannels = map[string]string{
			"1176204022399631381": "radio",
			"1118202946455351388": "music-1",
			"1118202975026937948": "music-2",
			"1118202999504904212": "music-3",
		}
	}

	if globalConfiguration.Environment == "staging" {
		// Dev dynamic channel creation button channels
		globalConfiguration.DynamicChannelCreateButtonIds = map[string]string{
			"1217251206624186481": "☕ | Dev Test Room (~Extra~)",
			"1217914805478887424": "🔒 | Dev Test Private Room (~Extra~)",
		}
	} else {
		// Production dynamic channel creation button channels
		globalConfiguration.DynamicChannelCreateButtonIds = map[string]string{
			"1171570400891785266": "☕ | Chill Room (~Extra~)",
			"1171589545473613886": "🔒 | Private Room (~Extra~)",
			"1171591013354197062": "🔮 | Spiritual Room (~Extra~)",
			"1171595498185035796": "🎵 | Music Room (~Extra~)",
			"1171599680568832023": "🎮 | Gaming (~Extra~)",
		}
	}

	if globalConfiguration.Environment == "staging" {
		// Dev default text channels
		globalConfiguration.DefaultInformationChannels = map[string]string{
			"1188135110042734613": "default",
			"1194451477192773773": "staff-rules",
			"1198686819928264784": "server-rules",
			"1205859615406030868": "legends",
		}
	} else {
		// Production default text channels
		globalConfiguration.DefaultInformationChannels = map[string]string{
			"1176277764001767464": "info-music",
			"1100486860058398770": "staff-rules",
			"1100142572141281460": "server-rules",
			"1100762035450544219": "legends",
		}
	}

}

// Load the available notification channels in the global map.
func LoadNotificationChannels() {

	for _, channelPairString := range globalConfiguration.NotificationChannelsPairs {

		isVoice, descriptor, channelId := getChannelValuesFromChannelPair(channelPairString)

		if descriptor == nil || channelId == nil {
			continue
		}

		globalConfiguration.NotificationChannels[*descriptor] = dataModels.Channel{
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
