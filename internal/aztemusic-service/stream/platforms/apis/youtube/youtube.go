package youtubeApi

import (
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/stream/platforms"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/download"
)

type YouTubeApi struct {
	ApiName string
}

// Downloads the mp3 file for the provided YT URL and returns data about it.
func (a YouTubeApi) GetMp3ForSongWithUrl(url string) platforms.DownloadedMp3Data {

	downloadTimestamp := time.Now().Unix()

	songName, mp3Data, err := download.DownloadMp3FromYoutube("mp3", url)
	if err != nil {
		log.Fatalf("Error occured while retrieving MP3 audio data for song: %s", err)
	}

	// Retrieve more song details to return to streamer
	// TODO

	return platforms.DownloadedMp3Data{
		Timestamp: downloadTimestamp,
		SongName:  songName,
		Artist:    "def",
		Duration:  "def",
		Mp3Data:   mp3Data,
	}
}
