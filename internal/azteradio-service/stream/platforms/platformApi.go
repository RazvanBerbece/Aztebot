package platforms

// An interface which provides basic functions to download and use MP3 audio data locally from various sources (local files, URLs, etc.).
// For example, a wrapper for YouTube downloads can be implemented with this interface.
type PlatformApi interface {
	GetMp3ForSongWithUrl(url string)
	// GetIdForSongName(songName string)
}
