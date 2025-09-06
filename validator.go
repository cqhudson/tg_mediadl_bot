package main

import (
    "log"

    "github.com/dlclark/regexp2"
)

func validateMessageContainsUrl(message string, regex string) bool {
	log.Printf("The message we are validating is --> %s", message)
	regexFormatted := regexp2.MustCompile(regex, 0)
	log.Printf("The compiled regex is --> %s", regexFormatted)

	matchFound, err := regexFormatted.MatchString(message)
	if err != nil {
		log.Printf("There was an error trying to find a match --> %s", err.Error())
		return false
	}

	return matchFound
}
