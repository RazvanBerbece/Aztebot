package streams

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/azteradio-service/stream/platforms"
	"github.com/bwmarrin/discordgo"
)

type Stream struct {
	PlatformName  string
	PlatformApi   platforms.PlatformApi
	TargetChannel *discordgo.VoiceConnection
	Strategy      string
	SourceFolder  string
	Queue         []string
}

func NewStream(platformName string,
	platformApi platforms.PlatformApi,
	targetChannel *discordgo.VoiceConnection) Stream {
	stream := Stream{
		PlatformName:  platformName,
		PlatformApi:   platformApi,
		TargetChannel: targetChannel,
	}
	stream.Queue = make([]string, 0)
	return stream
}

func (s *Stream) WithUrlsSourceFile() *Stream {
	s.Strategy = "URL-SRC-FOLDER"
	s.SourceFolder = "internal/azteradio-service/stream/sources/url"
	return s
}

func (s *Stream) PlayFromSource(ds *discordgo.Session) {

	switch s.Strategy {
	case "URL-SRC-FOLDER":
		fmt.Printf("Playing stream from source with strategy %s", s.Strategy)
		urls := s.getAllYoutubeUrls()
		for _, url := range urls {
			s.PlatformApi.GetMp3ForSongWithUrl(url)
			s.Queue = append(s.Queue, url)
		}
		// dgvoice.PlayAudioFile(s.TargetChannel, fmt.Sprintf("%s/%s", s.Folder, f.Name()), make(chan bool))
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
		readFile, err := os.Open(f.Name())
		if err != nil {
			log.Fatalf("Error occured while retrieving available YT URLs: %s", err)
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
