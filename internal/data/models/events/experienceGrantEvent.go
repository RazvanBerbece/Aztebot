package events

type ExperienceGrantEvent struct {
	UserId string
	Points float64
	Type   string
}
