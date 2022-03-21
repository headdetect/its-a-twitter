package store

import (
	"log"

	"github.com/headdetect/its-a-twitter/api/models"
)

var Sessions map[string]models.User = make(map[string]models.User) // [authToken] = user
var Users map[string]models.User = make(map[string]models.User)// [username] = user

var Tweets []models.Tweet = make([]models.Tweet, 1)
var Timelines map[string]models.Timeline = make(map[string]models.Timeline) // [username] = timeline

func GenerateTimelines() {
}

/*
	Restores the state from disk
*/
func Restore() {

}

/*
	Flushes the state to disk on a timed basis
*/
func StartFlusher() {
	log.Println("Flushing to DB")

	Restore()
}
