package main

import (
	"context"
	"log"
	"os" // needed to load env variables loaded from godotenv

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
    tu "github.com/mymmrac/telegoutil"
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

		const ytRegex string = "(?:https?://(?:(?:www\\.)?youtube\\.com/watch\\?v=|youtu\\.be/)[A-Za-z0-9_-]{11}(?:\\?si=[A-Za-z0-9_-]+)?)(?=\\s|$)"
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

            log.Printf("Attempting to download YouTube video")
			downloadedVideo, err = downloadYouTubeVideo(message, youtubeId)
            if err != nil {
                log.Printf("There was an error trying to download the video: %s", err.Error())
                continue
            }

            // let's send the downloaded video to the user now
            // func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
            sendVideoParams := telego.SendVideoParams {
                ChatID: tu.ID(update.Message.Chat.Id),
                Video: telego.InputFile { 
                    File: downloadedVideo,
                },
            }   
            sentMsg, err := bot.SendVideo(contect.Background(), *sendVideoParams)
            if err != nil {
                log.Printf("An error occurred while sending the video --> %s", err)
                continue
            }   

			continue
		}

	}
}

func checkForYouTubeLinks(message string, regex string) bool {
	return validateMessageContainsUrl(message, regex)
}
