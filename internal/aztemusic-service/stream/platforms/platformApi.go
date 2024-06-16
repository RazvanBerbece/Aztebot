package platforms

type DownloadedMp3Data struct {
	Timestamp int64
	SongName  string
	Artist    string
	Duration  string
	Mp3Data   []byte
}

// An interface which provides basic functions to download and use MP3 audio data locally from various sources (local files, URLs, etc.).
// For example, a wrapper for YouTube downloads can be implemented with this interface.
type PlatformApi interface {
	GetMp3ForSongWithUrl(url string) DownloadedMp3Data
	// GetIdForSongName(songName string)
}
