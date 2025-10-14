package main

import (
	"context"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/cqhudson/logger"
	_ "github.com/glebarez/go-sqlite"
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

		// TODO: implement an optional whitelist so the bot only responds to approved users
		/*
		something like this...
			if USING_WHITELIST && USER_ID IS NOT IN WHITELIST {
				continue
			}
		*/
		// Connect to SQLite db
		db, err := sql.Open("sqlite", "./db/tg_mediadl_bot.db")
		if err != nil {
			l.Printf("Failed to connect to SQLite database --> %s", err.Error())
		} else {
			l.Printf("Connected to SQLite database successfully.")
		}
		defer db.Close()

		// Get the version of SQLite
		var sqliteVersion string
		err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
		if err != nil {
			l.Printf("failed to get sqlite version --> %s", err.Error())
		} else {
			l.Printf("SQLite version --> %s", sqliteVersion)
		}

		l.Printf("Update:  %+v", update)

		// This regex is DISGUSTING and makes me sad
		//
		const ytRegex string = `https?://(?:www\.|m\.)?youtube\.com/(?:watch\?v=|shorts/)[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?|https?://youtu\.be/[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?`
		l.Printf("Checking for YouTube video using the following regex --> %s", ytRegex)

		containsYouTubeLink := checkForYouTubeLinks(update.Message.Text, ytRegex)

		if containsYouTubeLink {
			handleYouTubeVideo(&update, ytRegex, bot)
		}

		// This regex also disgusts me, just not as much. I need to determine a better way of checking link types
		//
		const xRegex string = `https?://(?:www\.)?(?:x|twitter)\.com/[a-zA-Z0-9_]+/status/\d+(?:/video/\d+)?(?:\?.*)?`
		l.Printf("Checking for X/Twitter video using the following regex --> %s", xRegex)

		containsXLink := checkForXLinks(update.Message.Text, xRegex)

		if containsXLink {
			handleXVideo(&update, xRegex, bot)
		}

	}
	//
	////
}

func checkForYouTubeLinks(message string, regex string) bool {
	return validateMessageContainsUrl(message, regex)
}

func checkForXLinks(message string, regex string) bool {
	return validateMessageContainsUrl(message, regex)
}

