package youtubeApi

import (
	"fmt"
	"log"
	"time"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/download"
)

type YouTubeApi struct {
	ApiName string
}

func (a YouTubeApi) GetMp3ForSongWithUrl(url string) {
	downloadTimestamp := time.Now().Unix()
	err := download.DownloadFile(fmt.Sprint(downloadTimestamp), "mp3", url)
	if err != nil {
		log.Fatalf("Error occured while retrieving MP3 audio data for song: %s", err)
	}
}
