package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/dlclark/regexp2"
)

func downloadYouTubeVideo(url string, youtubeId string) (*os.File, error) {

	binary, err := exec.LookPath("exec/yt-dlp")
	if err != nil {
		log.Printf("Failed to LookPath --> %s", err.Error())
	}
	log.Printf("Binary found --> %+v", binary)

	// Command line options for yt-dlp
	outputOption := fmt.Sprintf("download/yt/%s.%%(ext)s", youtubeId)
	compressionOptions := fmt.Sprintf("-S res:480")
	log.Printf("output options --> %s", outputOption)
	log.Printf("compression options --> %s", compressionOptions)

	args := []string{
		binary,
		"-o",
		outputOption,
		compressionOptions,
		"-f",
		"bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
		"--recode-video",
		"mp4",
		url,
	}

	env := os.Environ()

	err = syscall.Exec(binary, args, env)
	if err != nil {
		log.Printf("Error executing command --> %s", err.Error())
		return nil, err
	}

	filePath := filepath.Join("download/yt" + youtubeId + ".mp4")
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open downloaded file %s: %s", filePath, err.Error())
		return nil, err
	}
    log.Printf("file is %+v", file)

	return file, nil
}

func extractYouTubeId(url *regexp2.Match) (string, error) {

	log.Printf("Attempting to extract youtube ID from the following URL --> %s", url.Group.Capture.String())

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
	}

	// First let's grab the domain in the URL
	for i, letter := range fullUrl {
		// log.Printf("parsing youtu.be domain --> %s", string(letter))
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
			// 	log.Printf("parsing youtu.be link --> %s", string(fullUrl[i]))
			youtubeId += string(fullUrl[i])
		}
		return youtubeId, nil
	}

	if domain == "https://youtube.com" || domain == "http://youtube.com" || domain == "https://www.youtube.com" || domain == "http://www.youtube.com" {
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
			// implement later to support Live downloads
			// if temp := "live" {}
			// implement later to support Shorts downloads
			// if temp := "short" {}
		}
	}

	return "", errors.New("Unable to extract a YouTube ID")
}
