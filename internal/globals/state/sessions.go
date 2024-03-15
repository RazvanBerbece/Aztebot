package globalState

import "time"

var VoiceSessions = make(map[string]time.Time)
var StreamSessions = make(map[string]*time.Time)
var MusicSessions = make(map[string]map[string]*time.Time)
var DeafSessions = make(map[string]time.Time)

var ConfessionsToApprove = make(map[string]string)
