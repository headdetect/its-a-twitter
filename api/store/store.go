/*
	=== Notes about scalability ===

	This file exists to store the backend state of our application.
	In a normal application, this would need to persist in some fashion usually with
	a combination of a disk-persisted database (eg. PostgreSQL) and an in-memory database (eg. Redis)

	That would allow for a more generalized approach, allowing each service to handle
	it's own filtering, sorting, limiting. Additionally, the largest downside of a solution like this,
	is a lack of scalability due to a single file based database and the lack of ability to shard and scale
	those databases.

	In the interest of time, a memory-based storage with a tailor-made schema will
	suffice.

*/

package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/headdetect/its-a-twitter/api/utils"
	_ "github.com/mattn/go-sqlite3"
)

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
