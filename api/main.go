package main

import (
	"log"

	"github.com/joho/godotenv"
)


func main() {
	log.Println("Starting its-a-twitter API")

	log.Println("Loading .env")
	err := godotenv.Load()

	if err != nil {
    log.Fatal("Error loading .env file")
	}

	log.Println("Loading database")
	// db

	StartRouter()
}