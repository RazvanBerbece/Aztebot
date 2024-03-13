package events

type VoiceChannelCreateEvent struct {
	Name            string
	Private         bool
	ParentChannelId string // the channel ID which generated this event (assuming that new channel events are created only through the existing VC approach)
	Description     string
	ParentMemberId  string // who initiated the event
	ParentGuildId   string // the guild to create the new channel in
}
