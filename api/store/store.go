package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/headdetect/its-a-twitter/api/models"
)

/*
	This file exists to store the backend state of our application.
	In a normal application, this would need to persist in some fashion.
	If it came down to it, I'd most likely use a combination of postgres & redis.
	That would allow for a more generalized approach, allowing each service to handle
	it's own filtering, sorting, limiting. Additionally, the largest downside is a lack of
	ACID, meaning that the chance of loss of data is possible.

	However, in the interest of time, a memory-based storage with a tailor-made schema will
	suffice. Just don't plan to have this be production grade.
*/

var db *sql.DB

// Memory stores //
var Sessions map[string]*models.User = make(map[string]*models.User) // [authToken] = user
var Timelines map[string]models.Timeline = make(map[string]models.Timeline) // [username] = timeline

var Tweets []models.Tweet = make([]models.Tweet, 1)


func GetUserWithPassByUsername(username string) (*models.User, string, error) {
	var user models.User
	var hashedPassword string

	err := db.
		QueryRow(
			"select id, username, displayName, profilePicHash, password, createdAt from users where username = ? limit 1",
			username,
		).
		Scan(&user.Id, &user.Username, &user.DisplayName, &user.ProfilePicPath, &hashedPassword, &user.CreatedAt)

	return &user, hashedPassword, err
}

func GetUserById(id int) (*models.User, error) {
	var user models.User

	err := db.
		QueryRow(
			"select id, username, displayName, profilePicHash, createdAt from user where id = ? limit 1", 
			id,
		).
		Scan(&user.Id, &user.Username, &user.DisplayName, &user.ProfilePicPath, &user.CreatedAt)

	return &user, err
}

func LoadDatabase() {
	LoadDatabaseFromFile("./store/store.db", "./store/initial.sql")
}

func LoadDatabaseFromFile(databaseFile string, initialQueryFile string) {
	_, err := os.Stat(databaseFile); 
	existed := err == nil

	data, err := sql.Open("sqlite3", databaseFile)
	db = data // FIXME: Gotta be a better way to do this

	if err != nil {
		log.Fatal(err)
	}

	if !existed {
		initialQuery, err := os.ReadFile(initialQueryFile)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Seeding database")
		_, err = db.Exec(string(initialQuery))
	}

	if err != nil {
		log.Fatalf("%q\n", err)
	}
}