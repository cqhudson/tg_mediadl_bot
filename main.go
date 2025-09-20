package main

import (
	"context"
	"log"
	"fmt"
	"os" // needed to load env variables loaded from godotenv

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
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

	// Get updates from Telegram
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	// Loop through all updates that come in
	for update := range updates {
		log.Printf("Update:  %+v \n", update)

		message := update.Message.Text

		// This regex is DISGUSTING and makes me sad
		const ytRegex string = `https?://(?:www\.|m\.)?youtube\.com/watch\?v=[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?|https?://youtu\.be/[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?`
		containsYouTubeLink := checkForYouTubeLinks(message, ytRegex)
		if containsYouTubeLink {
			log.Print("The message contained a valid YouTube link. Attempting to download the YouTube video.")
			url, err := extractUrl(message, ytRegex)
			if err != nil {
				log.Printf("There was an issue extracting the YouTube URL --> %s", err.Error())
				continue
			}
			log.Printf("The extracted *Match object returned the following --> %+v", url)

			youtubeId, err := extractYouTubeId(url)
			if err != nil {
				log.Printf("There was an error trying to extract a YouTube ID from the given URL: %s", err.Error())
				continue
			}
			log.Printf("The extracted YouTube ID was --> %s", youtubeId)

			// If the video is already downloaded, we will just send that file
			log.Printf("Attempting to download YouTube video")
			downloadedVideo, err := downloadYouTubeVideo(message, youtubeId)
			if err != nil {
				log.Printf("There was an error trying to download the video: %s", err.Error())
				continue
			}

			// Let's send a message letting the user know that the video is being sent
			//
			message := fmt.Sprintf("Attempting to download %s, please be patient as larger videos may take some time to send", url)
			err = sendTelegramMessage(bot, &update, message, true)

			// let's send the downloaded video to the user now
			// func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
			sentMsg, err := bot.SendVideo(context.Background(), &telego.SendVideoParams{
					ChatID: tu.ID(update.Message.Chat.ID),
					Video: telego.InputFile{
					File: downloadedVideo,
				},
			})
			if err != nil {
				log.Printf("An error occurred while sending the video --> %s", err)
				continue
			}
			log.Printf("Sent Msg object contains the following --> %+v", sentMsg)

			continue
		}

	}
}

func checkForYouTubeLinks(message string, regex string) bool {
	return validateMessageContainsUrl(message, regex)
}

// Utility func for easily sending plain Telegram messages
func sendTelegramMessage(bot *telego.Bot, update *telego.Update, msg string, shouldLog bool) error {
	if shouldLog == true {
		username :=	update.Message.From.Username
		id := update.Message.From.ID
		log.Printf("Attempting to send the following message to user %s with ID %d --> %s", username, id, msg)
	}
		
	send, err := bot.SendMessage(context.Background(), &telego.SendMessageParams {
		ChatID: tu.ID(update.Message.Chat.ID),
		Text: msg,
		ProtectContent: true,
	})
	if err != nil {
		return fmt.Errorf("There was an issue sending a message: %s", err.Error())
	}

	if shouldLog == true {
		log.Printf("Message was sent successfully --> %+v", send)
	}
	return nil
}
