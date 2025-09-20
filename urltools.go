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

func extractUrl(message string, regex string) (*regexp2.Match, error) {
	regexFormatted := regexp2.MustCompile(regex, 0)

	match, err := regexFormatted.FindStringMatch(message)
	if err != nil {
		return nil, err
	} else if match == nil {
		errorMsg := "(extractUrl func) - match was empty or no match was able to be extracted."
		return nil, errors.New(errorMsg)
	}

	return match, nil
}
