package main

import (
    "log"
    "os/exec"
    "errors"

    "github.com/dlclark/regexp2"
)

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

func extractYouTubeId(url *regexp2.Match) string, error { 
    fullUrl := string(url)
    domain := "" 
    index := 0

    // First let's grab the domain in the URL 
    for i, letter := range fullUrl {
        domain += string(letter)
        if domain == "youtu.be" || domain == "youtube.com" {
            index = i
            break  
        }
    }
    if domain != "youtu.be" || domain != "youtube.com" {
        // Something really bad happened if you hit this block :(
        return "", errors.New("(extractYouTubeId func) - no valid YouTube domain could be extracted" 
    }


    youtubeId := ""

    if domain == "youtu.be" {
        // index+1 = the "/" char after the domain
        for i := index+1; i < len(fullUrl); i++ {
            if string(fullUrl[i]) == "?" || string(fullUrl[i]) == "\n" {
                return youtubeId, nil
            }
        }
           youtubeId += string(letter)
        }
    }

    return "", errors.New("Unable to extract a YouTube ID")

    // TODO - Need to come up with a solution for live videos and shorts. This will likely break otherwise
    //if domain == "youtube.com" {

    //}
}
