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

func extractYouTubeId(url *regexp2.Match) (string, error) { 

    log.Printf("Attempting to extract youtube ID from the following URL --> %s", url.Group.Capture.String())

    // url.Group.Capture.String() returns a string.
    fullUrl := url.Group.Capture.String()
    domain := "" 
    index := 0

    // First let's grab the domain in the URL 
    for i, letter := range fullUrl {
        log.Printf("parsing youtu.be domain --> %s", string(letter))
        domain += string(letter)
        if domain == "https://youtu.be" || domain == "http://youtu.be" || domain == "https://youtube.com" || domain == "http://youtube.com" {
            index = i
            break 
        }
    }
    if domain != "https://youtu.be" && domain != "http://youtu.be" && domain != "https://youtube.com" && domain != "http://youtube.com" {
        // Something really bad happened if you hit this block :(
        return "", errors.New("(extractYouTubeId func) - no valid YouTube domain could be extracted") 
    }


    youtubeId := ""

     
    log.Printf("Attempting to extract YT ID from domain %s", domain)
    if domain == "https://youtu.be" || domain == "http://youtu.be" {
        // index+1 = the "/" char after the domain
        for i := index+2; i < len(fullUrl); i++ {
            if string(fullUrl[i]) == "?" || string(fullUrl[i]) == "\n"  || string(fullUrl[i]) == " " {
                return youtubeId, nil
            }
            log.Printf("parsing youtu.be link --> %s", string(fullUrl[i]))
            youtubeId += string(fullUrl[i])
        }
        return youtubeId, nil
    }
    
    if domain == "https://youtube.com" || domain == "http://youtube.com" {
        temp := ""
        for i := index+2; i < len(fullUrl); i++ {
            if temp == "watch" {
                // this will skip the "?v=" chars 
                for j := i+3; j< len(fullUrl); j++ {
                    if string(fullUrl[j]) == "&" || string(fullUrl[j]) == "\n" || string(fullUrl[j]) == " " {
                        return youtubeId, nil
                    }
                    youtubeId += string(fullUrl[j]) 
                }

            }
            // implement later to support Live downloads
            // if temp := "live" {}
            // implement later to support Shorts downloads 
            // if temp := "short" {}
            log.Printf("parsing youtu.be link --> %s", string(fullUrl[i]))

            temp += string(fullUrl[i])
        }
    }

    return "", errors.New("Unable to extract a YouTube ID")
}
