package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
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

	l.Print("Your Telegram API token is %s", telegramToken) 
	//
	////

	//// Create the Telegram bot
	// 
	l.Print("Attempting to create a bot instance with Telegram")

	bot, err := telego.NewBot(telegramToken, telego.WithDefaultDebugLogger())
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

		// Message.Text contains the exact text the user sent the bot
		//
		message := update.Message.Text

		// Message.From.Username returns the username of the account who sent the msg
		//
		username := update.Message.From.Username


		// This regex is DISGUSTING and makes me sad
		//
		const ytRegex string = `https?://(?:www\.|m\.)?youtube\.com/watch\?v=[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?|https?://youtu\.be/[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?`
		l.Printf("Checking for YouTube video using the following regex --> %s", ytRegex)

		containsYouTubeLink := checkForYouTubeLinks(message, ytRegex, true)

		if containsYouTubeLink {

			// If the message contains a youtube link, attempt to extract it from the message
			//
			l.Print("The message contained a valid YouTube link. Attempting to extract it from the message.")

			url, err := extractUrl(message, ytRegex, true)

			if err != nil {
				l.Printf("There was an issue extracting the YouTube URL --> %s", err.Error())

				msg := "There was an issue extracting the youtube link from your message. Please ensure your URL is valid"
				err = sendTelegramMessage(bot, &update, msg)

				if err != nil {
					l.Printf("Failed to send message to %s --> %s", username, err.Error())
				} else {
					l.Printf("Successfully send %s the following Telegram message --> %s", username, msg)
				}

				continue
			}

			l.Printf("The extraction returned the following --> %+v", url)
			//
			////

			// If the URL is valid, we need to extract the YT ID from it
			//
			l.Print("Attempting to extract the ID from the YouTube link.")

			youtubeId, err := extractYouTubeId(url)

			if err != nil {
				l.Printf("There was an error trying to extract a YouTube ID from the given URL: %s", err.Error())
				
				msg := "There was an issue extracting the youtube ID from the link you sent. Please ensure your URL is valid."
				err = sendTelegramMessage(bot, &update, msg)

				if err != nil {
					l.Printf("Failed to send message to %s --> %s", username, err.Error())
				} else {
					l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
				}

				continue
			}

			l.Printf("The extracted YouTube ID was --> %s", youtubeId)
			//
			////

			// Attempt to download the video. If it is already downloaded, we will send the downloaded file. 
			//
			l.Printf("Attempting to download YouTube video")

			msg := "Attempting to download YouTube video. Please be patient as it can take a few minutes to finish the download."
			err = sendTelegramMessage(bot, &update, msg)

			if err != nil {
				l.Printf("Failed to send message to %s --> %s", username, err.Error())
			} else {
				l.Printf("Seccessfully sent %s the following Telegram message --> %s", username, msg)
			}

			downloadedVideo, err := downloadYouTubeVideo(message, youtubeId)

			if err != nil {
				l.Printf("There was an error trying to download the video: %s", err.Error())
				continue
			}

			l.Printf("The downloaded video was found --> %+v", downloadedVideo)
			//
			////

			// Let's send a message letting the user know that the video is being sent
			//
			l.Printf("Attempting to send a message to %s to let them know we are attempting to send them a file", username)

			msg = "Successfully downloaded the video. Attempting to send it to you. Please be patient as larger videos may take some time to send"
			err = sendTelegramMessage(bot, &update, msg)

			// TODO: Add error handling here

			l.Printf("Successfully sent the following message to %s --> %+v", username, msg)

			//
			////

			// let's send the downloaded video to the user now
			// func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
			//
			l.Printf("Attempting to send video file to %s", username)

			sentMsg, err := bot.SendVideo(context.Background(), &telego.SendVideoParams{
				ChatID: tu.ID(update.Message.Chat.ID),
				Video: telego.InputFile{
					File: downloadedVideo,
				},
			})

			if err != nil {
				l.Printf("An error occurred while sending the video --> %s", err)
				continue
			}

			l.Printf("Successfully sent the following message to %s --> %+v", username, sentMsg)
			//
			////

			continue
		}

	}
	//
	////
}

func checkForYouTubeLinks(message string, regex string, shouldLog bool) bool {
	return validateMessageContainsUrl(message, regex, shouldLog)
}

