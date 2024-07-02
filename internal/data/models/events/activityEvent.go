package events

type ActivityEvent struct {
	UserId string  // user that registered an activity
	Type   *string // MSG, REACT, VC, MUSIC, SLASH
	Value  *int64  // can be used to state how many units of activity to add
}
