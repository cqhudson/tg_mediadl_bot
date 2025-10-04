package main

import (
	"errors"
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/cqhudson/logger"
)

func validateMessageContainsUrl(message string, regex string) bool {

	l := logger.NewLogger(true, false, false)
	const loggingHeader string = "[validateMessageContainsUrl]"

	l.Printf("%s -- Attempting to validate the following message contains a URL --> %s", loggingHeader, message)

	regexFormatted, err := regexp2.Compile(regex, 0)
	if err != nil {
		l.Printf("%s -- ERROR compiling regex --> %s", loggingHeader, err.Error())
		return false
	}
	l.Printf("%s -- Compiled regex --> %s", loggingHeader, regexFormatted.String())

	matchFound, err := regexFormatted.MatchString(message)
	if err != nil {
		l.Printf("%s -- There was an error trying to find a match --> %s", loggingHeader, err.Error())
		return false
	}

	l.Printf("%s -- Matched the message \"%s\" against the compiled regex \"%s\"", loggingHeader, message, regex)
	return matchFound
}

func extractUrl(message string, regex string) (*regexp2.Match, error) {

	l := logger.NewLogger(true, false, false)
	const loggingHeader string = "[extractUrl]"

	l.Printf("%s -- attempting to extract URL from the following message --> %s", loggingHeader, message)
	l.Printf("%s -- matching against the following regex --> %s", loggingHeader, regex)

	regexFormatted := regexp2.MustCompile(regex, 0)

	match, err := regexFormatted.FindStringMatch(message)
	if err != nil {
		return nil, err
	} else if match == nil {
		errorMsg := fmt.Sprintf("%s -- match was empty or no match was able to be extracted.", loggingHeader)
		return nil, errors.New(errorMsg)
	}

	l.Printf("%s -- matched the following URL --> %+v", loggingHeader, match)
	return match, nil
}
