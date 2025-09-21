package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/cqhudson/logger"
)

func main() {
	
	// NewLogger(enablePrints, enablePanics, enableFatals)
	//
	logger := logger.NewLogger(true, false, true)

	const loggingPrefix string = "[main] - "
	logger.SetPrefix(loggingPrefix)

	logger.Print("logging is enabled")
	//
	////

	//// First load the .env file
	//
	logger.Print("Loading .env file")

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}
		
	logger.Print(".env file successfully loaded")
	//
	////

	//// Attempt to fetch Telegram API token after .env is loaded
	//
	logger.Print("Attempting to fetch Telegram API token from environment")

	telegramToken, telegramTokenExists := os.LookupEnv("TG_API_KEY")
	if telegramTokenExists == false {
		logger.Fatal("Error fetching TG_API_KEY environment variable")
	}

	logger.Print("Your Telegram API token is %s", telegramToken) 
	//
	////

	//// Create the Telegram bot
	// 
	logger.Print("Attempting to create a bot instance with Telegram")

	bot, err := telego.NewBot(telegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		logger.Fatalf("Failed to initialize bot: %s", err.Error())
	}

	logger.Printf("Successfully created bot instance: ID %d - Username %s", bot.ID(), bot.Username())
	//
	////

	//// Get updates from Telegram
	//
	logger.Print("Fetching updates from Telegram via long polling")

	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)
	//
	////

	//// Loop through all updates that come in
	//
	for update := range updates {

		logger.Printf("Update:  %+v", update)

		// Message.Text contains the exact text the user sent the bot
		//
		message := update.Message.Text

		// Message.From.Username returns the username of the account who sent the msg
		//
		username := update.Message.From.Username


		// This regex is DISGUSTING and makes me sad
		//
		const ytRegex string = `https?://(?:www\.|m\.)?youtube\.com/watch\?v=[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?|https?://youtu\.be/[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?`
		logger.Printf("Checking for YouTube video using the following regex --> %s", ytRegex)

		containsYouTubeLink := checkForYouTubeLinks(message, ytRegex, true)

		if containsYouTubeLink {

			// If the message contains a valid link, attempt to extract it from the message
			//
			logger.Print("The message contained a valid YouTube link. Attempting to extract it from the message.")

			url, err := extractUrl(message, ytRegex, true)

			if err != nil {
				logger.Printf("There was an issue extracting the YouTube URL --> %s", err.Error())
				// TODO: Send a Telegram message back to the user stating that the URL is invalid
				continue
			}

			logger.Printf("The extraction returned the following --> %+v", url)
			//
			////

			// If the URL is valid, we need to extract the YT ID from it
			//
			logger.Print("Attempting to extract the ID from the YouTube link.")

			youtubeId, err := extractYouTubeId(url)

			if err != nil {
				logger.Printf("There was an error trying to extract a YouTube ID from the given URL: %s", err.Error())
				// TODO: Send a Telegram message back to the user stating that the URL is invalid
				continue
			}

			logger.Printf("The extracted YouTube ID was --> %s", youtubeId)
			//
			////

			// Attempt to download the video. If it is already downloaded, we will send the downloaded file. 
			//
			logger.Printf("Attempting to download YouTube video")

			// TODO: Send a Telegram message to the user stating that we are attempting to dl the video

			downloadedVideo, err := downloadYouTubeVideo(message, youtubeId)

			if err != nil {
				logger.Printf("There was an error trying to download the video: %s", err.Error())
				continue
			}

			logger.Printf("The downloaded video was found --> %+v", downloadedVideo)
			//
			////

			// Let's send a message letting the user know that the video is being sent
			//
			logger.Printf("Attempting to send a message to %s to let them know we are attempting to send them a file", username)

			message := fmt.Sprintf("Attempting to send %s, please be patient as larger videos may take some time to send", url)
			err = sendTelegramMessage(bot, &update, message, true)

			// TODO: Add error handling here

			logger.Printf("Successfully sent the following message to %s --> %+v", username, message)

			//
			////

			// let's send the downloaded video to the user now
			// func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
			//
			logger.Printf("Attempting to send video file to %s", username)

			sentMsg, err := bot.SendVideo(context.Background(), &telego.SendVideoParams{
				ChatID: tu.ID(update.Message.Chat.ID),
				Video: telego.InputFile{
					File: downloadedVideo,
				},
			})

			if err != nil {
				logger.Printf("An error occurred while sending the video --> %s", err)
				continue
			}

			logger.Printf("Successfully sent the following message to %s --> %+v", username, sentMsg)
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

// Utility func for easily sending plain Telegram messages
func sendTelegramMessage(bot *telego.Bot, update *telego.Update, msg string, shouldLog bool) error {

	logger := logger.NewLogger(true, false, false)

	username := update.Message.From.Username
	id := update.Message.From.ID

	// Attempt to send the Telegram Message
	//
	logger.Printf("Attempting to send message to user %s with ID %d", username, id)

	send, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:         tu.ID(update.Message.Chat.ID),
		Text:           msg,
		ProtectContent: true,
	})

	if err != nil {
		return fmt.Errorf("There was an issue sending a message: %s", err.Error())
	}

	logger.Printf("Message was sent successfully --> %+v", send)
	//
	////

	return nil
}
