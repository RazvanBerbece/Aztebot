package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(filename string, fileType string, url string) (err error) {

	pathToLogfile := filepath.Join(".", fmt.Sprintf("downloads/%s/", fileType))
	err = os.MkdirAll(pathToLogfile, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create filepath for %s downloads: %s", fileType, err)
		return err
	}

	// Create the file to put the data in
	out, err := os.Create(fmt.Sprintf("./downloads/%s/%s.%s", fileType, filename, fileType))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while downloading data from %s: %s", url, resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
