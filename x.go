package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mymmrac/telego"
	"github.com/dlclark/regexp2"
	"github.com/cqhudson/logger"
)

func downloadVideoFromX(url string, xId string) (*os.File, error) {

	// Instantiate new Logger
	//
	l := logger.NewLogger(true, false, false)

	// First let's check for an existing video
	//
	existingFile, err := checkVideoAlreadyDownloaded("download/x", xId) // returns the path to the file

	if err != nil {
		l.Printf("There was an error trying to check for existing download --> %s", err.Error())
	} else {
		l.Print("Successfully checked for existing download")
	}

	// If an existing file was found, let's return a pointer to it instead of redownloading
	//
	if existingFile != "" {
		l.Printf("Existing file found at --> %s", existingFile)

		file, err := os.Open(existingFile)
		if err != nil {
			l.Printf("There was an error trying to fetch a pointer to the file --> %s", err.Error())
		} else {
			l.Printf("Got a pointer to the file --> %v", file)
			return file, nil
		}
	}
	
	// If a previously downloaded video doesn't exist, let's attempt to download
	//
	binary, err := exec.LookPath("exec/yt-dlp")
	if err != nil {
		l.Printf("Failed to LookPath --> %s", err.Error())
	}
	l.Printf("Binary found --> %+v", binary)

	// Command line options for yt-dlp
	outputFlag 		:= "-o"
	outputOptions 		:= fmt.Sprintf("download/x/%s.%%(ext)s", xId)

	compressionFlag		:= "-f"
	compressionOptions 	:= fmt.Sprintf("bestvideo[ext=mp4][height<=720]+bestaudio[ext=m4a][abr<=128]/best[ext=mp4][height<=720]")

	filesizeFlag 		:= "--max-filesize"
	filesizeOptions 	:= "49.9M"

	filetypeFlag 		:= "--remux-video"
	filetypeOptions 	:= "mp4"


	l.Printf("output options --> %s", outputOptions)
	l.Printf("compression options --> %s", compressionOptions)
	l.Printf("filesize options --> %s", filesizeOptions)
	l.Printf("filetype options --> %s", filetypeOptions)

	args := []string{
		outputFlag,		
		outputOptions,

		compressionFlag,
		compressionOptions,

		filesizeFlag,
		filesizeOptions,

		filetypeFlag,
		filetypeOptions,

		url,
	}

	// Attempt to download the video
	//
	l.Print("Attempting to download video")

	cmd := exec.Command(binary, args...)
	l.Printf("Running the following command --> %+v", cmd)
 
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		l.Printf("There was an error attempting to download the video --> %s", err.Error())
		l.Printf("Output from stdout/stderr -->\n%s", stdoutStderr)
		// TODO: Add logic here to determine if failure was caused by filesize.
		return nil, err
	}

	l.Printf("the output from running the command:\n%s", stdoutStderr)

	// Fetch a handle to the downloaded
	//
	filePath := filepath.Join("download/x/" + xId + ".mp4")
	file, err := os.Open(filePath)

	if err != nil {
		l.Printf("Failed to open downloaded file %s: %s", filePath, err.Error())
		return nil, err
	}

	return file, nil
}

func extractXId(url *regexp2.Match) (string, error) {
	
	l := logger.NewLogger(true, false, false)

	fullUrl := url.Group.Capture.String()
	domain := ""
	index := 0
	validDomains := map[string]bool{
		"https://x.com":		true,
		"http://x.com":			true,
		"https://www.x.com":		true,
		"http://www.x.com":		true,
		"https://twitter.com":		true,
		"http://twitter.com":		true,
		"https://www.twitter.com":	true,
	}

	// First let's grab the domain in the URL
	for i, letter := range fullUrl {
		domain += string(letter)
		if validDomains[domain] == true {
			index = i
			break
		}
	}
	l.Printf("The domain parsed was %s", domain)
	if validDomains[domain] != true {
		// Something bad happened if you hit this block :/
		return "", errors.New("(extractXId func) - no valid X/Twitter domain could be extracted")
	}

	xId := ""
	temp := ""
	substringLength := len("shorts/")

	if validDomains[domain] {
		// If the parsing is specific per domain, rewrite this by referring to youtube.go
		for i := index+1; i < len(fullUrl); i++ {
			temp += string(fullUrl[i])
			
			// This is getting the substring for the last 7 chars of the temp variable, searching for "status/"
			if len(temp) > substringLength {
				if string(temp[len(temp)-substringLength:len(temp)]) == "status/" {
					for j := i+1; j < len(fullUrl); j++ {
						// next chars are the video ID
						if string(fullUrl[j]) != "\n" && string(fullUrl[j]) != "/" && string(fullUrl[j]) != "" && string(fullUrl[j]) != " " && string(fullUrl[j]) != "?" {
							xId += string(fullUrl[j])
						} else {
							return xId, nil
						}
					}
				}
			}
		}	
	}
	return "", errors.New("Failed to parse out an X ID from the URL")
}

func handleXVideo(update *telego.Update, xRegex string, bot *telego.Bot) {
	
	l := logger.NewLogger(true, false, false)

	// Message.Text contains the exact text the user sent the bot
	//
	message := update.Message.Text

	// Message.From.Username returns the username of the account who sent the msg
	//
	username := update.Message.From.Username


	// If the message contains a youtube link, attempt to extract it from the message
	//
	l.Print("The message contained a valid X link. Attempting to extract it from the message.")
	url, err := extractUrl(message, xRegex)
	if err != nil {
		l.Printf("There was an issue extracting the X URL --> %s", err.Error())
		msg := "There was an issue extracting the X/Twitter link from your message. Please ensure your URL is valid"
		_ = sendTelegramMessage(bot, update, msg)
	}
	l.Printf("The extraction returned the following --> %+v", url)
	//
	//

	// If the URL is valid, we need to extract the YT ID from it
	//
	l.Print("Attempting to extract the ID from the X link.")
	xId, err := extractXId(url)
	if err != nil {
		l.Printf("There was an error trying to extract a X ID from the given URL: %s", err.Error())

		msg := "There was an issue extracting the X/Twitter ID from the link you sent. Please ensure your URL is valid."
		_ = sendTelegramMessage(bot, update, msg)
	}
	l.Printf("The extracted X ID was --> %s", xId)
	//
	//

	// Attempt to download the video. If it is already downloaded, we will send the downloaded file. 
	//
	msg := "Attempting to download X video. Please be patient as it can take a few minutes to finish the download."
	err = sendTelegramMessage(bot, update, msg)
	if err != nil {
		l.Printf("Failed to send message to %s --> %s", username, err.Error())
	} else {
		l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
	}
	downloadedVideo, err := downloadVideoFromX(message, xId)
	if err != nil {
		l.Printf("There was an error trying to download the video: %s", err.Error())

		msg := "Failed to download video. I appologize for the inconvenience. Please wait a little bit and try again later."
		_ = sendTelegramMessage(bot, update, msg)
	}
	l.Printf("The downloaded video was found --> %+v", downloadedVideo)
	//
	//

	// Let's send a message letting the user know that the video is being sent
	//
	msg = "Successfully downloaded the video. Attempting to send it to you. Please be patient as larger videos may take some time to send"
	err = sendTelegramMessage(bot, update, msg)
	if err != nil {
		l.Printf("Failed to send message to %s --> %s", username, err.Error())

		msg := "Failed to send video due to a network error. I appologize for the inconvenience. Please try again later."
		_ = sendTelegramMessage(bot, update, msg)
	}
	l.Printf("Successfully sent the following message to %s --> %+v", username, msg)
	//
	//

	// let's send the downloaded video to the user now
	// func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
	//
	l.Printf("Attempting to send video file to %s", username)
	err = SendTelegramVideo(bot, update, downloadedVideo)	
	if err != nil {
		l.Printf("An error occurred while sending the video --> %s", err.Error())
	}
	l.Printf("Successfully sent the video to %s", username)
	//
	//
}
