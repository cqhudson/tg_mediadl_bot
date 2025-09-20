// cleaner logging in your go program
//

package main

import (
	"log"
)

// Logger wraps logging functionality behind a toggle
//
type Logger struct {
	enabled bool	
}

// Create a new Logger instance 
//
func NewLogger(enabled bool) *Logger {
	return &Logger {enabled: enabled}
}

// Wrap log.Printf
//
func (logger *Logger) Printf(format string, args ...interface{}) {
	if logger.enabled == true {
		log.Printf(format, args...)
	}	
}

// Wrap log.Print
//
func (logger *Logger) Print(args ...interface{}) {
	if logger.enabled == true {
		log.Print(args)	
	}
}

