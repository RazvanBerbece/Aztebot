package member

import (
	"fmt"

	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/bwmarrin/discordgo"
)

// Checks against the local maps if a user has an active voice session.
func MemberHasActiveVoiceSession(uid string) bool {

	status := 0

	if _, ok := globalState.VoiceSessions[uid]; ok {
		status += 1
	}

	if _, ok := globalState.MusicSessions[uid]; ok {
		status += 1
	}

	if _, ok := globalState.StreamSessions[uid]; ok {
		status += 1
	}

	return status == 3

}

func GetUserVoiceChannel(s *discordgo.Session, guildID, userID string) (string, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return "", fmt.Errorf("error retrieving guild: %v", err)
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", fmt.Errorf("user is not in a voice channel")
}
