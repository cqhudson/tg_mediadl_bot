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

func downloadYouTubeVideo(url string, youtubeId string) (*os.File, error) {

	// Instantiate new logger
	//
	l := logger.NewLogger(true, false, false)

	// First let's check for an existing video
	//
	l.Printf("Checking for an existing YouTube video download for %s", youtubeId)

	existingDownload, err := checkVideoAlreadyDownloaded("download/yt", youtubeId)

	if err != nil {
		l.Printf("There was an error trying to check for existing download --> %s", err.Error())
	} else {
		l.Print("Successfully checked for existing download")
	}

	// If an existing download was found, let's return a pointer to it instead of redownloading
	//
	if existingDownload != "" {
		l.Printf("Existing file found at --> %s", existingDownload)

		file, err := os.Open(existingDownload)
		if err != nil {
			l.Printf("There was an error trying to fetch a pointer to the file --> %s", err.Error())
		} else {
			l.Printf("Got a pointer to the file --> %+v", file)
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
	outputFlag    := "-o"
	outputOptions := fmt.Sprintf("download/yt/%s.%%(ext)s", youtubeId)

	compressionFlag    := "-f"
	compressionOptions := fmt.Sprint("bestvideo[ext=mp4][height<=720]+bestaudio[ext=m4a][abr<=128]/best[ext=mp4][height<=720]")

	filesizeFlag    := "--max-filesize"
	filesizeOptions := "49.9M"

	filetypeFlag    := "--remux-video"
	filetypeOptions := "mp4"
	

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
	l.Print("attempting to fetch a handle to the downloaded file")

	filePath := filepath.Join("download/yt/" + youtubeId + ".mp4")
	file, err := os.Open(filePath)

	if err != nil {
		l.Printf("Failed to open downloaded file %s: %s", filePath, err.Error())
		return nil, err
	}

	fileInfo, _ := file.Stat()

	// If file is over 49.5MB, we shouldn't send it (Telegram Bot API Limitation)
	//
	if fileInfo.Size()

	// Max filesize supported
	// 49.9 MB (in bytes)
	const var MAX_FILESIZE = 49999999
	if fileInfo.Size() > MAX_FILESIZE {
		l.Printf("Video filesize is too large to send --> %d", fileInfo.Size())
		// TODO: Initiate a file cleanup here (Delete the video)
		return nil, errors.New("Filesize too large to to send to user.")
	} else {
		l.Printf("filename found is %s and the filesize is %d bytes", fileInfo.Name(), fileInfo.Size())
		return file, nil
	}
}

func extractYouTubeId(url *regexp2.Match) (string, error) {

	l := logger.NewLogger(true, false, false)

	// url.Group.Capture.String() returns a string.
	fullUrl := url.Group.Capture.String()
	domain := ""
	index := 0
	validDomains := map[string]bool{
		"https://youtu.be":        true,
		"http://youtu.be":         true,
		"https://youtube.com":     true,
		"http://youtube.com":      true,
		"https://www.youtube.com": true,
		"http://www.youtube.com":  true,
		"https://m.youtube.com":   true,
		"http://m.youtube.com":    true,
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
	if validDomains[domain] == false {
		// Something really bad happened if you hit this block :(
		return "", errors.New("(extractYouTubeId func) - no valid YouTube domain could be extracted")
	}

	youtubeId := ""

	if domain == "https://youtu.be" || domain == "http://youtu.be" {
		// index+1 = the "/" char after the domain
		for i := index + 2; i < len(fullUrl); i++ {
			if string(fullUrl[i]) == "?" || string(fullUrl[i]) == "\n" || string(fullUrl[i]) == " " {
				return youtubeId, nil
			}
			youtubeId += string(fullUrl[i])
		}
		return youtubeId, nil
	}

	if domain == "https://youtube.com" || domain == "http://youtube.com" || domain == "https://www.youtube.com" || domain == "http://www.youtube.com" || domain == "https://m.youtube.com" || domain == "http://m.youtube.com" {

		temp := ""
		for i := index + 2; i < len(fullUrl); i++ {
			l.Printf("temp var == %s", temp)
			if temp == "watch" {
				// this will skip the "?v=" chars
				offset := 3
				for j := i + offset; j < len(fullUrl); j++ {
					if string(fullUrl[j]) == "&" || string(fullUrl[j]) == "\n" || string(fullUrl[j]) == " " {
						return youtubeId, nil
					}
					youtubeId += string(fullUrl[j])
				}
				return youtubeId, nil

			}

			// implement later to support Shorts downloads
			// if temp := "short" {}
			if temp == "shorts" {
				// this will skip the initial "/" after shorts in the url
				offset := 1
				for j := i + offset; j < len(fullUrl); j++ {
					if string(fullUrl[j]) == "?" || string(fullUrl[j]) == "\n" || string(fullUrl[j]) == " " {
						return youtubeId, nil
					}
					youtubeId += string(fullUrl[j])
				}
				return youtubeId, nil
			}

			temp += string(fullUrl[i])
		}
	}

	return "", errors.New("Unable to extract a YouTube ID")
}

// This is called in main.go to facilitate downloading a YouTube video
func handleYouTubeVideo(update *telego.Update, ytRegex string, bot *telego.Bot) {

	l := logger.NewLogger(true, false, false)

	// Message.Text contains the exact text the user sent the bot
	//
	message := update.Message.Text

	// Message.From.Username returns the username of the account who sent the msg
	//
	username := update.Message.From.Username


	// If the message contains a youtube link, attempt to extract it from the message
	//
	l.Print("The message contained a valid YouTube link. Attempting to extract it from the message.")

	url, err := extractUrl(message, ytRegex)

	if err != nil {
		l.Printf("There was an issue extracting the YouTube URL --> %s", err.Error())

		msg := "There was an issue extracting the youtube link from your message. Please ensure your URL is valid"
		err = sendTelegramMessage(bot, update, msg)

		if err != nil {
			l.Printf("Failed to send message to %s --> %s", username, err.Error())
		} else {
			l.Printf("Successfully send %s the following Telegram message --> %s", username, msg)
		}
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
		err = sendTelegramMessage(bot, update, msg)

		if err != nil {
			l.Printf("Failed to send message to %s --> %s", username, err.Error())
		} else {
			l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
		}
	}

	l.Printf("The extracted YouTube ID was --> %s", youtubeId)
	//
	////

	// Attempt to download the video. If it is already downloaded, we will send the downloaded file. 
	//
	l.Printf("Attempting to download YouTube video")

	msg := "Attempting to download YouTube video. Please be patient as it can take a few minutes to finish the download."
	err = sendTelegramMessage(bot, update, msg)

	if err != nil {
		l.Printf("Failed to send message to %s --> %s", username, err.Error())
	} else {
		l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
	}

	downloadedVideo, err := downloadYouTubeVideo(message, youtubeId)

	if err != nil {
		l.Printf("There was an error trying to download the video: %s", err.Error())

		msg := "Failed to download video. I appologize for the inconvenience. Please wait a little bit and try again later."
		err = sendTelegramMessage(bot, update, msg)
		 
		if err != nil {
			l.Printf("Failed to send message to %s --> %s", username, err.Error())
		} else { 
			l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
		}
	}

	l.Printf("The downloaded video was found --> %+v", downloadedVideo)
	//
	////

	// Let's send a message letting the user know that the video is being sent
	//
	l.Printf("Attempting to send a message to %s to let them know we are attempting to send them a file", username)

	msg = "Successfully downloaded the video. Attempting to send it to you. Please be patient as larger videos may take some time to send"
	err = sendTelegramMessage(bot, update, msg)

	if err != nil {
		l.Printf("Failed to send message to %s --> %s", username, err.Error())
	} else { 
		l.Printf("Successfully sent %s the following Telegram message --> %s", username, msg)
	}

	l.Printf("Successfully sent the following message to %s --> %+v", username, msg)

	//
	////

	// let's send the downloaded video to the user now
	// func (b *Bot) SendVideo(ctx context.Context, params *SendVideoParams) (*Message, error)
	//
	l.Printf("Attempting to send video file to %s", username)
	
	err = SendTelegramVideo(bot, update, downloadedVideo)	

	if err != nil {
		l.Printf("An error occurred while sending the video --> %s", err)
	}

	l.Printf("Successfully sent the video to %s", username)
	//
	////
}
