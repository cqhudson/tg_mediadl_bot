package main

import (
	"log"
	"os" // needed to load env variables loaded from godotenv

	"github.com/joho/godotenv"
)

func main() {

	// First load the .env file
	log.Print("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Print(".env file successfully loaded")

	// Attempt to fetch Telegram API token after .env is loaded
	telegramToken, telegramTokenExists := os.LookupEnv("TG_API_KEY")
	if telegramTokenExists == false {
		log.Fatal("Error fetching TG_API_KEY environment variable")
	}
	log.Printf("Your Telegram API token is %s", telegramToken) // This line is for local debugging and should be disabled when running in prod
}
