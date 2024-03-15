package member

import "github.com/RazvanBerbece/Aztebot/internal/globals"

// Checks against the global maps if a user has an active voice session.
func MemberHasActiveVoiceSession(uid string) bool {

	status := 0

	if _, ok := globals.VoiceSessions[uid]; ok {
		status += 1
	}

	if _, ok := globals.MusicSessions[uid]; ok {
		status += 1
	}

	if _, ok := globals.StreamSessions[uid]; ok {
		status += 1
	}

	return status == 3

}
