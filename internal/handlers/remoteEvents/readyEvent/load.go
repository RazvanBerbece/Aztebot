package readyEvent

import (
	"strings"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/domain"
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
		//
		// DEVELOPMENT STATIC CHANNELS
		//

		// AFK CHANNELS
		globalConfiguration.AfkChannels = map[string]string{
			"1176284686297874522": "afk",
		}

		// MUSIC CHANNELS
		globalConfiguration.MusicChannels = map[string]string{
			"1173790229258326106": "radio",
		}

		// DYNAMIC CHANNEL CREATION BUTTON CHANNELS
		globalConfiguration.DynamicChannelCreateButtonIds = map[string]string{
			"1217251206624186481": "☕ | Dev Test Room (~Extra~)",
			"1217914805478887424": "🔒 | Dev Test Private Room (~Extra~)",
		}

		// DEFAULT INFO TEXT CHANNELS
		globalConfiguration.DefaultInformationChannels = map[string]string{
			"1194451477192773773": "staff-rules",
			"1198686819928264784": "server-rules",
			"1205859615406030868": "legends",
			"1188135110042734613": "socials",
		}

		// CHANNELS WHICH ARE EXCEPTED FROM AUTOMATIC MESSAGE DELETIONS
		globalConfiguration.DeleteExceptedChannels = map[string]string{
			"1213272204326998056": "jail",
		}
	} else {
		//
		// PRODUCTION STATIC CHANNELS
		//

		// AFK CHANNELS
		globalConfiguration.AfkChannels = map[string]string{
			"1212508073101627412": "afk",
		}

		// MUSIC CHANNELS
		globalConfiguration.MusicChannels = map[string]string{
			"1176204022399631381": "radio",
			"1118202946455351388": "music-1",
			"1118202975026937948": "music-2",
			"1118202999504904212": "music-3",
		}

		// DYNAMIC CHANNEL CREATION BUTTON CHANNELS
		globalConfiguration.DynamicChannelCreateButtonIds = map[string]string{
			"1171570400891785266": "☕ | Chill Room (~Extra~)",
			"1171589545473613886": "🔒 | Private Room (~Extra~)",
			"1171591013354197062": "🔮 | Spiritual Room (~Extra~)",
			"1171595498185035796": "🎵 | Music Room (~Extra~)",
			"1171599680568832023": "🎮 | Gaming (~Extra~)",
		}

		// DEFAULT INFO TEXT CHANNELS
		globalConfiguration.DefaultInformationChannels = map[string]string{
			"1176277764001767464": "info-music",
			"1100486860058398770": "staff-rules",
			"1100142572141281460": "server-rules",
			"1219417926482788443": "legends",
			"1168327391937044500": "socials",
		}

		// CHANNELS WHICH ARE EXCEPTED FROM AUTOMATIC MESSAGE DELETIONS
		globalConfiguration.DeleteExceptedChannels = map[string]string{
			"1214338314010886234": "🏛・celula",
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

		globalConfiguration.NotificationChannels[*descriptor] = domain.Channel{
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
