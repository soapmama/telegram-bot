package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TOKEN")

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Call method getMe
	botUser, _ := bot.GetMe()
	fmt.Printf("Bot User: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	for update := range updates {
		if update.Message != nil || update.ChannelPost != nil {
			// Retrieve chat ID
			chatID := update.ChannelPost.Chat.ID

			// Call method sendMessage.
			// Send a message to sender with the same text (echo bot).
			// (https://core.telegram.org/bots/api#sendmessage)
			sentMessage, _ := bot.SendMessage(
				tu.Message(
					tu.ID(chatID),
					"Мыльная папа советует: 🧼 Покупайте наше мыло 🧼",
				),
			)

			fmt.Printf("Sent Message: %v\n", sentMessage)
		}
	}
}
