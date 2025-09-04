package main

import (
	"context"
	"log"
	"os" // needed to load env variables loaded from godotenv
	"os/exec"
	//"strings"

	"github.com/dlclark/regexp2" // more feature-rich regex package based on the .NET regex engine
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

		ytRegexPattern := "(?:https?://(?:(?:www\\.)?youtube\\.com/watch\\?v=|youtu\\.be/)[A-Za-z0-9_-]{11}(?:\\?si=[A-Za-z0-9_-]+)?)(?=\\s|$)"

		ytInfo := newYouTubeInformation(ytRegexPattern)

		// Let's try and download a YouTube video depending on the input
		isValidYTVideo := validateUrlIsYouTube(message, ytInfo)
		log.Printf("Is the message sent a valid YouTube URL? --> %s", isValidYTVideo)
		if isValidYTVideo == true {
			//downloadYouTubeVideo(message)
			log.Print("this is a valid youtube video link")
		}
	}

}

func validateUrlIsYouTube(message string, ytInfo *YouTubeInformation) bool {
	// These are both the same video, we must be able to search both
	// Example valid URL: https://youtu.be/pOcg-AdC2Y8?si=RC5W58mNkSLO2omb
	// Example valid URL: https://www.youtube.com/watch?v=pOcg-AdC2Y8
	// TODO -> Also need to validate if the link is missing the "http://" or "https://" substr.

	// HERE IS AN EXAMPLE REGEX PATTERN
	/*
	     *    \b(?:https?://(?:(?:www\\.)?youtube\\.com/watch\\?v=|youtu\\.be/)[A-Za-z0-9_-]{11}(?:\\?si=[A-Za-z0-9_-]+)?)(?=\\s|$)
	     *    The above regex is compiled below to:
	     *    (?:https?://(?:(?:www\.)?youtube\.com/watch\?v=|youtu\.be/)[A-Za-z0-9_-]{11}(?:\?si=[A-Za-z0-9_-]+)?)(?=\s|$)
		 *
		 * Matches the first YouTube video URL (e.g., https://youtube.com/watch?v=VIDEO_ID or https://youtu.be/VIDEO_ID,
		 *    with optional www. and http://, and optional ?si= query for youtu.be)
		 *    with an 11-character video ID (letters, digits, hyphens, underscores),
		 *    stopping at a space, newline, or end of string.
	*/

	log.Printf("The message we are validating is --> %s", message)

	// regexp2.MustCompile parses a regex and returns, if successful, a pointer to a Regexp object to match against text
	regex := regexp2.MustCompile(ytInfo.Regex, 0)
	log.Printf("The compiled regex is --> %s", regex.String())

	matchFound, err := regex.MatchString(message)
	if err != nil {
		log.Printf("There was an error trying to find a match --> %s", err.Error())
		return false
	}

	return matchFound
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
