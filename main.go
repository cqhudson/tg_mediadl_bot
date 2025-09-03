package main

import (
	"log"
	"os" // needed to load env variables loaded from godotenv

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
)

func main() {

	// First load the .env file
	log.Print("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}
	log.Print(".env file successfully loaded")

	// Attempt to fetch Telegram API token after .env is loaded
	telegramToken, telegramTokenExists := os.LookupEnv("TG_API_KEY")
	if telegramTokenExists == false {
		log.Fatal("Error fetching TG_API_KEY environment variable")
	}
	log.Printf("Your Telegram API token is %s", telegramToken) // This line is for local debugging and should be disabled when running in prod

	// Create the Telegram bot and enable debugging info
	// (debugging info should only be used during local development)
	bot, err := telego.NewBot(telegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Failed to initialize bot: %s", err.Error())
	}
	log.Printf("Bot user: %+v \n", bot)

}
