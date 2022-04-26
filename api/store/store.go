package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/headdetect/its-a-twitter/api/utils"
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

func LoadDatabase(overwrite bool) {
	appEnv, _ := utils.GetStringOrDefault("APP_ENV", "development")
	storePath, _ := utils.GetStringOrDefault("STORE_PATH", "./store")
	dbFilePath := fmt.Sprintf("%s/%s.db", storePath, appEnv)

	if overwrite {
		os.Remove(dbFilePath)
	}

	existed := LoadDatabaseFromFile(dbFilePath, "rcw")

	if !existed {
		SeedDatabase(
			fmt.Sprintf("%s/schema.sql", storePath),
			fmt.Sprintf("%s/seeds/%s.sql", storePath, appEnv),
		)
	}
}

func LoadDatabaseFromFile(databaseFile string, openMode string) bool {
	_, err := os.Stat(databaseFile)
	existed := err == nil

	data, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=%s", databaseFile, openMode))
	DB = data

	DB.SetMaxOpenConns(1)

	if err != nil {
		log.Fatal(err)
	}

	return existed
}

func SeedDatabase(files ...string) {
	for _, file := range files {
		sql, err := os.ReadFile(file)

		if err != nil {
			log.Println("Could not find seed file: ", file)
			continue
		}

		_, err = DB.Exec(string(sql))

		if err != nil {
			log.Fatalf("%q\n", err)
		}
	}
}
