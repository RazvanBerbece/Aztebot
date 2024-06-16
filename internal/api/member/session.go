package member

import globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"

// Checks against the global maps if a user has an active voice session.
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
