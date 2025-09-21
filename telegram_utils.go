package main

import (
	"fmt"
	"context"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/cqhudson/logger"
)

// Utility func for easily sending plain Telegram messages
func sendTelegramMessage(bot *telego.Bot, update *telego.Update, msg string, shouldLog bool) error {

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
