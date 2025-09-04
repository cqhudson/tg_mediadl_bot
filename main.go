package main

import (
	"context"
	"log"
	"os" // needed to load env variables loaded from godotenv
	"os/exec"
	"strings"

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

	// Call method getMe (https://core.telegram.org/bots/api#getme)
	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		log.Fatalf("Failed to get bot information: %s", err.Error())
	}
	log.Printf("Bot user: %+v \n", botUser)

	// Get updates from Telegram
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	// Loop through all updates that come in
	for update := range updates {
		log.Printf("Update:  %+v \n", update)

		message := update.Message.Text

		// Let's try and download a YouTube video depending on the input
		isValidYTVideo := validateUrlIsYouTube(message)
		log.Printf("Is the message sent a valid YouTube URL? --> %b", isValidYTVideo)
		if isValidYTVideo == true {
			downloadYouTubeVideo(message)
		}
	}

}

func validateUrlIsYouTube(url string) bool {

	// Example valid URL: https://youtu.be/nvUTNX0FDPA

	// (TODO) Example flow for this function:
	// 1 - Check if first few chars are "https://" or "http://"
	// 2 - Check if the following substr is either "youtu.be" or "youtube.com"
	// 3 - Check if there is a space or newline at the end of the URL (We only want the URl, and nothing else)
	// 4 - Check if the link is a valid URL

	// 1 - Check if url contains http or https
	beginning := ""
	containsHTTP := false
	for _, letter := range url {
		beginning += string(letter)
		if beginning == "http://" || beginning == "https://" {
			containsHTTP = true
			break
		}
	}
	if containsHTTP == false {
		log.Printf("%s does not start with \"http://\" or \"http://\"")
		return false
	}


    // 2 - (TODO) rework this to check if next characters are a valid YouTube domain
    // instead of blindly checking for the existence of the domain.
	isValid := false
	validYouTubeDomains := []string{
		"youtu.be",
		"youtube.com",
	}

	log.Printf("Checking if message is a valid YouTube URL. Message --> %s", url)

	for _, domain := range validYouTubeDomains {
		log.Printf("domain: %s", domain)
		if (strings.Contains(url, domain)) == true {
			isValid = true
			break
		}
	}

	return isValid
}

func downloadYouTubeVideo(url string) error {

	// Command line options for yt-dlp
	// Example command: .\executables\yt-dlp.exe -o ".\downloads\YT\%(autonumber)06d.%(ext)s" https://youtu.be/n8-wN0lc5qk?si=cD1KaaffXHWjn0jq
	// Example: save video as ".\downloads\YT\000004.webm"
	outputOption := "-o \".\\downloads\\YT\\%(autonumber)06d.%(ext)s\""

	cmd := exec.Command("exec/yt-dlp.exe", outputOption, url)
	err := cmd.Run()
	if err != nil {
		log.Printf("Unable to execute command: %s", err.Error())
	}

	log.Printf("The command run was --> %s", cmd)
	return nil
}
