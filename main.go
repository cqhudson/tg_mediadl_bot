package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/cqhudson/logger"
)

func main() {
	
	// NewLogger(enablePrints, enablePanics, enableFatals)
	//
	l := logger.NewLogger(true, false, true)

	const loggingPrefix string = "[main] - "
	l.SetPrefix(loggingPrefix)
	l.SetFlags(logger.Lshortfile | logger.Ltime | logger.Lmicroseconds | logger.Lmsgprefix)

	l.Print("logging is enabled")
	//
	////

	//// First load the .env file
	//
	l.Print("Loading .env file")

	err := godotenv.Load()
	if err != nil {
		l.Fatalf("Error loading .env file: %s", err.Error())
	}

	l.Print(".env file successfully loaded")
	//
	////

	//// Attempt to fetch Telegram API token after .env is loaded
	//
	l.Print("Attempting to fetch Telegram API token from environment")

	telegramToken, telegramTokenExists := os.LookupEnv("TG_API_KEY")
	if telegramTokenExists == false {
		l.Fatal("Error fetching TG_API_KEY environment variable")
	}

	l.Printf("Your Telegram API token is %s", telegramToken) 
	//
	////

	//// Create the Telegram bot
	// 
	l.Print("Attempting to create a bot instance with Telegram")

	bot, err := telego.NewBot(telegramToken, telego.WithDiscardLogger())
	if err != nil {
		l.Fatalf("Failed to initialize bot: %s", err.Error())
	}

	l.Printf("Successfully created bot instance: ID %d - Username %s", bot.ID(), bot.Username())
	//
	////

	//// Get updates from Telegram
	//
	l.Print("Fetching updates from Telegram via long polling")

	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)
	//
	////

	//// Loop through all updates that come in
	//
	for update := range updates {

		l.Printf("Update:  %+v", update)

		// This regex is DISGUSTING and makes me sad
		//
		const ytRegex string = `https?://(?:www\.|m\.)?youtube\.com/(?:watch\?v=|shorts/)[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?|https?://youtu\.be/[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?`
		l.Printf("Checking for YouTube video using the following regex --> %s", ytRegex)

		containsYouTubeLink := checkForYouTubeLinks(update.Message.Text, ytRegex, true)

		if containsYouTubeLink {
			handleYouTubeVideo(&update, ytRegex, bot);
		}

	}
	//
	////
}

func checkForYouTubeLinks(message string, regex string, shouldLog bool) bool {
	return validateMessageContainsUrl(message, regex, shouldLog)
}

