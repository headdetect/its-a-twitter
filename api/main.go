package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/headdetect/its-a-twitter/api/controller"
	"github.com/headdetect/its-a-twitter/api/store"
	"github.com/headdetect/its-a-twitter/api/utils"

	"github.com/joho/godotenv"
)

// [Scaling]
// This is problem unique to this setup. Since we have an in-memory session store
// we have to save it on close of the application so we don't lose all of
// our users' sessions
func saveSessions() {
	var encodedBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encodedBuffer)

	if err := encoder.Encode(controller.Sessions); err != nil {
		log.Fatalln("Unable to encode session")
	}

	os.WriteFile("sessions.dat", encodedBuffer.Bytes(), 0600)
}

func restoreSessions() {
	buffer, err := os.ReadFile("sessions.dat")

	if err != nil {
		log.Println("No sessions file found. Starting new session")
	}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	decoder.Decode(&controller.Sessions)
}

func main() {
	log.Println("Starting its-a-twitter API")

	log.Println("Loading .env")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	assetsPath, _ := utils.GetStringOrDefault("MEDIA_PATH", "./assets/media")
	newPath := filepath.Join(".", assetsPath)
	err = os.MkdirAll(newPath, os.ModePerm)

	if err != nil {
		log.Fatal("Error making assets directory")
		log.Fatal(err)
	}

	log.Println("Loading database")
	store.LoadDatabase(false)

	log.Println("Restoring sessions")
	restoreSessions()

	// Start router in a different thread //
	go controller.StartRouter()

	// Listening for sigterm/sigints so we can dump
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	log.Println() // So it comes on the next line
	log.Println("Saving sessions")
	saveSessions()

	log.Println("Good Bye!")
}
