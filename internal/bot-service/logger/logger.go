package logger

import (
	"log"
	"os"
	"path/filepath"
)

func LogHandlerCall(handlerName string, filename string) {

	// if filename not provided, use default
	if filename == "" {
		filename = "handler_calls.log"
	}

	// If the file and/or path doesn't exist, create, or append to the file
	pathToLogfile := filepath.Join(".", "logs/")
	err := os.MkdirAll(pathToLogfile, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create filepath for Handler calls log file. Err: ", err)
	}
	file, err := os.OpenFile("logs/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Could not open log file for Handler calls. Err: ", err)
	}
	defer file.Close()

	log.SetOutput(file)

	log.Printf("%s", handlerName)

	// Set logging output back to stderr
	log.SetOutput(os.Stderr)

}
