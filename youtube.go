package main

import (
    "log"
    "os/exec"
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
