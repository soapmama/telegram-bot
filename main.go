package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type App struct {
	bot    *telego.Bot
	config *Config
}

func (app *App) handleChannelPost(update *telego.Update) {
	if update.Message != nil || update.ChannelPost != nil {
		// Retrieve chat ID
		chatID := update.ChannelPost.Chat.ID

		// Call method sendMessage.
		// Send a message to sender with the same text (echo bot).
		// (https://core.telegram.org/bots/api#sendmessage)
		sentMessage, _ := app.bot.SendMessage(
			tu.Message(
				tu.ID(chatID),
				"Мыльная папа советует: 🧼 Покупайте наше мыло 🧼",
			),
		)

		fmt.Printf("Sent Message: %v\n", sentMessage)
	}
}

func main() {
	config := newConfig()

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(config.Token, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := &App{
		bot:    bot,
		config: config,
	}

	// Call method getMe
	botUser, _ := bot.GetMe()
	fmt.Printf("Bot User: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	for update := range updates {
		app.handleChannelPost(&update)
	}
}
