package download

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func DownloadMp3FromYoutube(fileType string, url string) (songName string, mp3Data []byte, err error) {

	outputPath := fmt.Sprintf("downloads/%s/%s.mp3", fileType, "%(title)s")

	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "mp3", "--output", outputPath, url)

	var stdout bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	errCmd := cmd.Run()
	if errCmd != nil {
		return "", nil, fmt.Errorf("error executing command: %v", errCmd)
	}

	// Get song name from output filepath
	fileName := filepath.Base(outputPath)
	songName = fileName[:len(fileName)-len(filepath.Ext(fileName))]

	return songName, stdout.Bytes(), nil

}

func DownloadFile(filename string, fileType string, url string) (mp3Data io.ReadCloser, err error) {

	pathToDownloadsFolder := filepath.Join(".", fmt.Sprintf("downloads/%s/", fileType))
	err = os.MkdirAll(pathToDownloadsFolder, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create filepath for %s downloads: %s", fileType, err)
		return nil, err
	}

	// Create the file to put the data in
	out, err := os.Create(fmt.Sprintf("./downloads/%s/%s.%s", fileType, filename, fileType))
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while downloading data from %s: %s", url, resp.Status)
	}

	// Writer the body to file
	data := resp.Body
	_, err = io.Copy(out, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
