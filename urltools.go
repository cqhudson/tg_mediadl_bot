package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/dlclark/regexp2"
)

func validateMessageContainsUrl(message string, regex string, shouldLog bool) bool {
	const loggingHeader string = "[validateMessageContainsUrl]"
	if shouldLog == true {
		log.Printf("%s -- Attempting to validate the following message contains a URL --> %s", loggingHeader, message)
	}

	regexFormatted, err := regexp2.Compile(regex, 0)
	if err != nil {
		log.Printf("%s -- ERROR compiling regex --> %s", loggingHeader, err.Error())
		return false
	}
	if shouldLog == true {
		log.Printf("%s -- Compiled regex --> %s", loggingHeader, regexFormatted.String())
	}

	matchFound, err := regexFormatted.MatchString(message)
	if err != nil {
		log.Printf("%s -- There was an error trying to find a match --> %s", loggingHeader, err.Error())
		return false
	}

	if shouldLog == true {
		log.Printf("%s -- Matched the message \"%s\" against the compiled regex \"%s\"", loggingHeader, message, regex)
	}
	return matchFound
}

func extractUrl(message string, regex string, shouldLog bool) (*regexp2.Match, error) {
	const loggingHeader string = "[extractUrl]"
	if shouldLog == true {
		log.Printf("%s -- attempting to extract URL from the following message --> %s", loggingHeader, message)
		log.Printf("%s -- matching against the following regex --> %s", loggingHeader, regex)
	}

	regexFormatted := regexp2.MustCompile(regex, 0)

	match, err := regexFormatted.FindStringMatch(message)
	if err != nil {
		return nil, err
	} else if match == nil {
		errorMsg := fmt.Sprintf("%s -- match was empty or no match was able to be extracted.", loggingHeader)
		return nil, errors.New(errorMsg)
	}

	if shouldLog == true {
		log.Printf("%s -- matched the following URL --> %+v", loggingHeader, match)
	}
	return match, nil
}
