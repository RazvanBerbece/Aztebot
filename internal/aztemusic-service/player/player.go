package player

type Song struct {
	FileName      string
	SongName      string
	TotalDuration string
	Artist        string
}

type Player struct {
	Queue        []Song // FIFO
	CurrentSong  Song
	PreviousSong Song
	NextSong     Song
	IsPaused     bool
}

func NewPlayer() *Player {
	player := Player{
		Queue:    make([]Song, 0),
		IsPaused: true,
	}
	return &player
}

// Plays the first song on the queue.
func (p *Player) Play() {
	// Play song - send audio data to Discord etc.
	p.CurrentSong = p.GetFrontOfQueue()

	// Mutate queue
}

// Adds a song to the front of the player queue (updates other relevant player fields).
func (p *Player) AddSongToQueue(song Song) {
	p.Queue = append(p.Queue, song)
}

// Retrieves the front song of the player queue (i.e song to be played).
func (p *Player) GetFrontOfQueue() Song {
	return p.Queue[0]
}

// Removes the front song of the player queue (updates other relevant player fields).
func (p *Player) RemoveFrontOfQueue() {
	p.Queue = p.Queue[1:]
}
