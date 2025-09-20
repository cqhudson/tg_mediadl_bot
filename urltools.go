package main

import (
	"errors"
	"log"

	"github.com/dlclark/regexp2"
)

func validateMessageContainsUrl(message string, regex string) bool {
	log.Printf("The message we are validating is --> %s", message)
	regexFormatted := regexp2.MustCompile(regex, 0)
	log.Printf("The compiled regex is --> %s", regexFormatted)

	matchFound, err := regexFormatted.MatchString(message)
	if err != nil {
		log.Printf("(MatchString) There was an error trying to find a match --> %s", err.Error())
		return false
	}

	return matchFound
}

func extractUrl(message string, regex string, shouldLog bool) (*regexp2.Match, error) {
	if shouldLog == true {
		log.Printf("[extractUrl] -- attempting to extract URL from the following message --> %s", message)
		log.Printf("[extractUrl] -- matching against the following regex --> %s", regex)
	}

	regexFormatted := regexp2.MustCompile(regex, 0)

	match, err := regexFormatted.FindStringMatch(message)
	if err != nil {
		return nil, err
	} else if match == nil {
		errorMsg := "[extractUrl] - match was empty or no match was able to be extracted."
		return nil, errors.New(errorMsg)
	}

	if shouldLog == true {
		log.Printf("[extractUrl] -- matched the following URL --> %+v", match)
	}
	return match, nil
} 
