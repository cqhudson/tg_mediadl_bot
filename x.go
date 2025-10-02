package main

import (
	"github.com/cqhudson/logger"
)

func downloadVideoFromX(url string, xid string) (*os.File, error) {
	// TODO: Implement support for downloading videos from Xitter

	// Example links to download from:
	// https://x.com/Solopopsss/status/1973363399234052232
	// https://x.com/theo/status/1973210960522559746
	// https://x.com/theo/status/1973167911419412985
	// https://x.com/theo/status/1973167911419412985/video/1

	// Seems like /video/1 is a redirect to the post itself.


}

func extractXId(url *regexp2.Match) (string, error) {
	
	l := logger.NewLogger(true, false, false)

	fullUrl := url.Group.Capture.String()
	domain := ""
	index := 0
	validDomains := map[string]bool{
		"https://x.com":	true,
		"http://x.com":		true,
		"https://twitter.com":	true,
		"http://twitter.com":	true,
	}

	// First let's grab the domain in the URL
	for i, letter := range fullURL {
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

	if validDomains[domain] {
		// If the parsing is specific per domain, rewrite this by referring to youtube.go
		for i := index+1; i < len(fullUrl); i++ {
			temp := string(fullUrl[i])

			// This is getting the substring for the last 7 chars of the temp variable, searching for "status/"
			if len(temp) > 5 && string(temp[len(temp)-6:len(temp)]) == "status/" {
				for j := i++; j < len(fullUrl); j++ {
					// next chars are the video ID
					if string(fullUrl[i]) != "\n" || string(fullUrl[i]) != "/" || string(fullUrl[i]) != "" || string(fullUrl[i]) != " " {
						xId += string(fullUrl[i])
					}
				}
				return xId, nil
			}
		}	
	}
	return "", errors.New("Failed to parse out an X ID from the URL")

}
