package main

import (
	"fmt"
	"context"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/cqhudson/logger"
)

// Utility func for easily sending plain Telegram messages
func sendTelegramMessage(bot *telego.Bot, update *telego.Update, msg string) error {

	l := logger.NewLogger(true, false, false)

	username := update.Message.From.Username
	id := update.Message.From.ID

	// Attempt to send the Telegram Message
	//
	l.Printf("Attempting to send message to user %s with ID %d", username, id)

	send, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:         tu.ID(update.Message.Chat.ID),
		Text:           msg,
		ProtectContent: true,
	})

	if err != nil {
		return fmt.Errorf("There was an issue sending a message: %s", err.Error())
	}

	l.Printf("Message was sent successfully --> %+v", send)
	//
	////

	return nil
}

func SendTelegramVideo(bot *telego.Bot, update *telego.Update, video *os.File) error {

	l := logger.NewLogger(true, false, false)

	username := update.Message.From.Username
	id := update.Message.From.ID

	// Attempt to send the Telegram Message
	//
	l.Printf("Attempting to send message to user %s with ID %d", username, id)
	
	send, err := bot.SendVideo(context.Background(), &telego.SendVideoParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Video: telego.InputFile {
			File: video,
		},
	})

	if err != nil {
		return fmt.Errorf("There was an issue sending a video: %s", err.Error())
	}

	l.Printf("Message was sent successfully --> %+v", send)
	//
	////

	return nil 
}
