package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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

var DB *sql.DB

func LoadDatabase() {
	LoadDatabaseFromFile("./store/store.db", "./store/initial.sql", "rcw")
}

func LoadDatabaseFromFile(databaseFile string, initialQueryFile string, openMode string) {
	_, err := os.Stat(databaseFile); 
	existed := err == nil

	data, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=%s", databaseFile, openMode))
	DB = data // FIXME: Gotta be a better way to do this

	DB.SetMaxOpenConns(1)

	if err != nil {
		log.Fatal(err)
	}

	if !existed {
		initialQuery, err := os.ReadFile(initialQueryFile)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Seeding database")
		_, err = DB.Exec(string(initialQuery))

		if err != nil {
			log.Fatalf("%q\n", err)
		}
	}

	if err != nil {
		log.Fatalf("%q\n", err)
	}
}