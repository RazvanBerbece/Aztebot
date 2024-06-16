package streams

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/player"
	"github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/stream/platforms"
	"github.com/bwmarrin/discordgo"
)

type Stream struct {
	PlatformName  string
	PlatformApi   platforms.PlatformApi
	TargetChannel *discordgo.VoiceConnection
	Strategy      string
	SourceFolder  string
	Player        *player.Player
}

func NewStream(platformName string,
	platformApi platforms.PlatformApi,
	targetChannel *discordgo.VoiceConnection) Stream {
	stream := Stream{
		PlatformName:  platformName,
		PlatformApi:   platformApi,
		TargetChannel: targetChannel,
	}

	stream.Player = player.NewPlayer()

	go stream.startPlayerQueueConsumer()

	return stream
}

func (s *Stream) WithUrlsSourceFile() *Stream {
	s.Strategy = "URL-SRC-FOLDER"
	s.SourceFolder = "internal/aztemusic-service/stream/sources/url"
	return s
}

func (s *Stream) PlayFromLocalSource(ds *discordgo.Session) {

	switch s.Strategy {
	case "URL-SRC-FOLDER":
		fmt.Printf("Playing stream from source with strategy %s", s.Strategy)
		urls := s.getAllYoutubeUrls()
		for id, url := range urls {
			go s.downloadMp3AndAddToQueue(url)
			if id > 5 {
				break
			}
		}
	case "MP3-SRC-FOLDER":
		log.Fatal("Not implemented")
	}

	time.Sleep(250 * time.Millisecond)

}

func (s *Stream) getAllYoutubeUrls() []string {

	var urls []string

	files, _ := os.ReadDir(s.SourceFolder)
	for _, f := range files {
		// Read URLs from file
		urlFileFullPath := fmt.Sprintf("%s/%s", s.SourceFolder, f.Name())
		readFile, err := os.Open(urlFileFullPath)
		if err != nil {
			log.Fatalf("Error occured while retrieving available YT URLs from (%s): %s", urlFileFullPath, err)
			return nil
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		// Append to output array
		for fileScanner.Scan() {
			urls = append(urls, fileScanner.Text())
		}
		readFile.Close()
	}

	return urls
}

func (s *Stream) downloadMp3AndAddToQueue(url string) {
	song := s.PlatformApi.GetMp3ForSongWithUrl(url)
	s.Player.AddSongToQueue(player.Song{
		FileName:      fmt.Sprintf("downloads/mp3/%d", song.Timestamp),
		SongName:      song.SongName,
		TotalDuration: song.Duration,
		Artist:        song.Artist,
	})
}

func (s *Stream) startPlayerQueueConsumer() {
	for {
		if len(s.Player.Queue) > 0 {
			// There are songs on the queue, so consume queue elements and play them
			fmt.Printf("Playing %s\n", s.Player.GetFrontOfQueue().FileName)
			// s.Player.Play()
		} else {
			return
		}
		time.Sleep(time.Second)
	}
}
