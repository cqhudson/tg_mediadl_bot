package main

import (
	"errors"
    "fmt"
	"log"
	"os/exec"
    "path/filepath"

	"github.com/dlclark/regexp2"
)

func downloadYouTubeVideo(url string, youtubeId string) error {

    // Resolve absolute path to yt-dlp executable
    absYtDlpPath, err := filepath.Abs("./exec/yt-dlp.exe")
    if err != nil {
        log.Printf("Failed to resolve yt-dlp path %s", err.Error())
        return err
    }   
    log.Printf("The following path was resolved --> %s", absYtDlpPath)  

    // Resolve absolute path to downloads directory 
    downloadPath := filepath.Join(".", "downloads", "YT")
    absDownloadPath, err := filepath.Abs(downloadPath)
    if err != nil {
        log.Printf("Failed to resolve download path %s", err.Error())
        return err
    }  
    log.Printf("The following path was resolved --> %s", absDownloadPath)

       

	// Command line options for yt-dlp
	outputOption := fmt.Sprintf("-o %s\\%s.%%(ext)s", absDownloadPath, youtubeId)
    compressionOptions := fmt.Sprintf("-S \"res:480\"")
    log.Printf("output options --> %s", outputOption)
    log.Printf("compression options --> %s", compressionOptions)

    args := fmt.Sprintf("%s %s", outputOption, compressionOptions)
    log.Printf("The args for this command is %s", args)

	cmd := exec.Command(absYtDlpPath, args, url)
    stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Unable to execute command: %+v", err)
	}
    log.Printf("%s/n", stdoutStderr)

	log.Printf("The command run was --> %s", cmd.String())
	return nil
}

func extractYouTubeId(url *regexp2.Match) (string, error) {

	log.Printf("Attempting to extract youtube ID from the following URL --> %s", url.Group.Capture.String())

	// url.Group.Capture.String() returns a string.
	fullUrl := url.Group.Capture.String()
	domain := ""
	index := 0
    validDomains := map[string]bool {
        "https://youtu.be": true,
        "http://youtu.be": true,
        "https://youtube.com": true,
        "http://youtube.com": true,
        "https://www.youtube.com": true,
        "http://www.youtube.com": true,
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

    if (domain == "https://youtu.be" || domain == "http://youtu.be") {
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
