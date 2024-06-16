package streams

import (
	"fmt"
	"os"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/azteradio-service/stream/platforms"
	"github.com/bwmarrin/discordgo"
)

type Stream struct {
	PlatformName  string
	PlatformApi   platforms.PlatformApi
	TargetChannel *discordgo.VoiceConnection
	Folder        string
}

func (s *Stream) DownloadSongs() {
	s.PlatformApi.GetMp3ForSongWithUrl("https://www.youtube.com/watch?v=-5YdR7GcUGU")
}

func (s *Stream) Play(ds *discordgo.Session) {

	// Start loop and attempt to play all files in the given folder
	fmt.Println("Reading Folder: ", s.Folder)
	files, _ := os.ReadDir(s.Folder)
	for _, f := range files {
		fmt.Println("Play:", f.Name())
		// dgvoice.PlayAudioFile(s.TargetChannel, fmt.Sprintf("%s/%s", s.Folder, f.Name()), make(chan bool))
	}

	time.Sleep(250 * time.Millisecond)

}
