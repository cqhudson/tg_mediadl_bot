package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/cqhudson/logger"
)

func downloadYouTubeVideo(url string, youtubeId string) (*os.File, error) {

	// Instantiate new logger
	//
	l := logger.NewLogger(true, false, false)
	//

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

	l.Printf("filename found is %s and the filesize is %d bytes", fileInfo.Name(), fileInfo.Size())

	return file, nil
}

func extractYouTubeId(url *regexp2.Match) (string, error) {

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
	log.Printf("The domain parsed was %s", domain)
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
			log.Printf("temp var == %s", temp)
			if temp == "watch" {
				// this will skip the "?v=" chars
				for j := i + 3; j < len(fullUrl); j++ {
					if string(fullUrl[j]) == "&" || string(fullUrl[j]) == "\n" || string(fullUrl[j]) == " " {
						return youtubeId, nil
					}
					youtubeId += string(fullUrl[j])
				}
				return youtubeId, nil

			}
			temp += string(fullUrl[i])
			// implement later to support Live downloads
			// if temp := "live" {}
			// implement later to support Shorts downloads
			// if temp := "short" {}
		}
	}

	return "", errors.New("Unable to extract a YouTube ID")
}

// if the video is already downloaded, then we can send the existing video
func checkVideoAlreadyDownloaded(dir string, filename string) (string, error) {
	var filePath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			// Get filename without extension since we don't necessarily know
			// what the ext will be (mp4, webm, etc)
			nameWithoutExt := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			if nameWithoutExt == filename {
				filePath = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("Failed to find existing file with base name of %s", filename)
	}
	return filePath, nil
}
